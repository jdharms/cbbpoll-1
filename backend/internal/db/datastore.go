package db

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"

	"github.com/r-cbb/cbbpoll/internal/errors"
	"github.com/r-cbb/cbbpoll/pkg"
)

// Eventually rename DatastoreClient to DatastoreClient and abstract out an interface type DatastoreClient
type DatastoreClient struct {
	client *datastore.Client
}

type idStruct struct {
	ID int64
}

func (i *idStruct) Load(property []datastore.Property) error {
	var ok, foundId bool
	for _, v := range property {
		if v.Name == "ID" {
			i.ID, ok = v.Value.(int64)
			if !ok {
				return fmt.Errorf("error loading ID property")
			}
			foundId = true
		}
	}
	if !foundId {
		return fmt.Errorf("no ID property on load")
	}
	return nil
}

func (i idStruct) Save() ([]datastore.Property, error) {
	return nil, fmt.Errorf("Should never save an idStruct to storage")
}

func NewDatastoreClient(projectId string) (*DatastoreClient, error) {
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, projectId)
	if err != nil {
		return nil, err
	}

	// Verify that we can communicate and authenticate with the datastore service.
	t, err := client.NewTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not connect: %v", err)
	}
	if err := t.Rollback(); err != nil {
		return nil, fmt.Errorf("datastoredb: could not connect: %v", err)
	}

	return &DatastoreClient{client: client}, nil
}

func (db *DatastoreClient) nextID(kind string) (id int64, err error) {
	ctx := context.Background()
	q := datastore.NewQuery(kind).Order("-ID")
	var ids []idStruct

	_, err = db.client.GetAll(ctx, q, &ids)
	if err != nil {
		return
	}

	if len(ids) == 0 {
		id = 1
	} else {
		id = ids[0].ID + 1
	}

	return
}

func (db *DatastoreClient) AddTeam(team pkg.Team) (id int64, err error) {
	ctx := context.Background()

	newId, err := db.nextID("Team")
	if err != nil {
		fmt.Printf("error finding next ID: %v", err.Error())
		return 0, err
	}
	team.ID = newId
	k := datastore.IDKey("Team", newId, nil)

	tx, err := db.client.NewTransaction(ctx)
	if err != nil {
		return 0, fmt.Errorf("datastoredb: could not create transaction: %v", err)
	}

	var tmp pkg.Team

	// Perform a Get or Put to ensure atomicity
	err = tx.Get(k, &tmp)
	if err == nil || err != datastore.ErrNoSuchEntity{
		_ = tx.Rollback()
		return 0, fmt.Errorf("concurrency error adding Team")
	}

	pk, err := tx.Put(k, &team)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("datastoredb: could not put team entity: %v", err)
	}

	c, err := tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("datastoredb: error committing transaction: %v", err)
	}

	k = c.Key(pk)
	if k.ID != newId {
		panic("keys don't match")
	}

	return newId, nil
}

func (db *DatastoreClient) GetTeam(id int64) (team pkg.Team, err error) {
	const op errors.Op = "datastore.GetTeam"
	ctx := context.Background()

	k := datastore.IDKey("Team", id, nil)
	err = db.client.Get(ctx, k, &team)

	if err == datastore.ErrNoSuchEntity {
		err = errors.E(errors.KindNotFound, op, err)
	} else if err != nil {
		err = errors.E(op, err)
	}

	return
}