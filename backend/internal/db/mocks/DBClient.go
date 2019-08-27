// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import models "github.com/r-cbb/cbbpoll/internal/models"

// DBClient is an autogenerated mock type for the DBClient type
type DBClient struct {
	mock.Mock
}

// AddTeam provides a mock function with given fields: newTeam
func (_m *DBClient) AddTeam(newTeam models.Team) (models.Team, error) {
	ret := _m.Called(newTeam)

	var r0 models.Team
	if rf, ok := ret.Get(0).(func(models.Team) models.Team); ok {
		r0 = rf(newTeam)
	} else {
		r0 = ret.Get(0).(models.Team)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Team) error); ok {
		r1 = rf(newTeam)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddUser provides a mock function with given fields: newUser
func (_m *DBClient) AddUser(newUser models.User) (models.User, error) {
	ret := _m.Called(newUser)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(models.User) models.User); ok {
		r0 = rf(newUser)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.User) error); ok {
		r1 = rf(newUser)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTeam provides a mock function with given fields: id
func (_m *DBClient) GetTeam(id int64) (models.Team, error) {
	ret := _m.Called(id)

	var r0 models.Team
	if rf, ok := ret.Get(0).(func(int64) models.Team); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(models.Team)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTeams provides a mock function with given fields:
func (_m *DBClient) GetTeams() ([]models.Team, error) {
	ret := _m.Called()

	var r0 []models.Team
	if rf, ok := ret.Get(0).(func() []models.Team); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Team)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: name
func (_m *DBClient) GetUser(name string) (models.User, error) {
	ret := _m.Called(name)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(string) models.User); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
