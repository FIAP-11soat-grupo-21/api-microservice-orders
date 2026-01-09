package factories

import (
	"testing"

	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
	"microservice/internal/interfaces"
)

type mockOrderDataSource struct{}

func (m *mockOrderDataSource) Create(order daos.OrderDAO) error {
	return nil
}

func (m *mockOrderDataSource) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	return []daos.OrderDAO{}, nil
}

func (m *mockOrderDataSource) FindByID(id string) (daos.OrderDAO, error) {
	return daos.OrderDAO{}, nil
}

func (m *mockOrderDataSource) Update(order daos.OrderDAO) error {
	return nil
}

func (m *mockOrderDataSource) Delete(id string) error {
	return nil
}

type mockOrderStatusDataSource struct{}

func (m *mockOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	return daos.OrderStatusDAO{}, nil
}

func (m *mockOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	return []daos.OrderStatusDAO{}, nil
}

func TestSetNewOrderDataSource(t *testing.T) {
	mockDS := &mockOrderDataSource{}

	SetNewOrderDataSource(func() interfaces.IOrderDataSource {
		return mockDS
	})

	ds := NewOrderDataSource()
	if ds == nil {
		t.Error("NewOrderDataSource() returned nil after SetNewOrderDataSource")
	}
	if ds != mockDS {
		t.Error("NewOrderDataSource() did not return the mock after SetNewOrderDataSource")
	}
}

func TestSetNewOrderDataSource_Nil(t *testing.T) {
	mockDS := &mockOrderDataSource{}
	SetNewOrderDataSource(func() interfaces.IOrderDataSource {
		return mockDS
	})

	SetNewOrderDataSource(nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("SetNewOrderDataSource(nil) restored default behavior")
		}
	}()
}

func TestSetNewOrderStatusDataSource(t *testing.T) {
	mockDS := &mockOrderStatusDataSource{}

	SetNewOrderStatusDataSource(func() interfaces.IOrderStatusDataSource {
		return mockDS
	})

	ds := NewOrderStatusDataSource()
	if ds == nil {
		t.Error("NewOrderStatusDataSource() returned nil after SetNewOrderStatusDataSource")
	}
	if ds != mockDS {
		t.Error("NewOrderStatusDataSource() did not return the mock after SetNewOrderStatusDataSource")
	}
}

func TestSetNewOrderStatusDataSource_Nil(t *testing.T) {
	mockDS := &mockOrderStatusDataSource{}
	SetNewOrderStatusDataSource(func() interfaces.IOrderStatusDataSource {
		return mockDS
	})

	SetNewOrderStatusDataSource(nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("SetNewOrderStatusDataSource(nil) restored default behavior")
		}
	}()
}

func TestNewOrderDataSource_WithMock(t *testing.T) {
	mockDS := &mockOrderDataSource{}

	SetNewOrderDataSource(func() interfaces.IOrderDataSource {
		return mockDS
	})
	defer SetNewOrderDataSource(nil)

	ds := NewOrderDataSource()
	if ds != mockDS {
		t.Error("NewOrderDataSource() did not return the mock")
	}
}

func TestNewOrderStatusDataSource_WithMock(t *testing.T) {
	mockDS := &mockOrderStatusDataSource{}

	SetNewOrderStatusDataSource(func() interfaces.IOrderStatusDataSource {
		return mockDS
	})
	defer SetNewOrderStatusDataSource(nil)

	ds := NewOrderStatusDataSource()
	if ds != mockDS {
		t.Error("NewOrderStatusDataSource() did not return the mock")
	}
}
func TestNewOrderDataSource_Default(t *testing.T) {
	SetNewOrderDataSource(nil)
	
	ds := NewOrderDataSource()
	if ds == nil {
		t.Error("Expected data source to be created")
	}
}

func TestNewOrderStatusDataSource_Default(t *testing.T) {
	SetNewOrderStatusDataSource(nil)
	
	ds := NewOrderStatusDataSource()
	if ds == nil {
		t.Error("Expected data source to be created")
	}
}

func TestSetNewOrderDataSource_CustomFunction(t *testing.T) {
	originalFn := newOrderDataSource
	defer func() {
		newOrderDataSource = originalFn
	}()
	
	customCalled := false
	customFn := func() interfaces.IOrderDataSource {
		customCalled = true
		return &mockOrderDataSourceCustom{}
	}
	
	SetNewOrderDataSource(customFn)
	
	ds := NewOrderDataSource()
	
	if !customCalled {
		t.Error("Expected custom function to be called")
	}
	
	if _, ok := ds.(*mockOrderDataSourceCustom); !ok {
		t.Error("Expected custom data source type")
	}
}

func TestSetNewOrderStatusDataSource_CustomFunction(t *testing.T) {
	originalFn := newOrderStatusDataSource
	defer func() {
		newOrderStatusDataSource = originalFn
	}()
	
	customCalled := false
	customFn := func() interfaces.IOrderStatusDataSource {
		customCalled = true
		return &mockOrderStatusDataSourceCustom{}
	}
	
	SetNewOrderStatusDataSource(customFn)
	
	ds := NewOrderStatusDataSource()
	
	if !customCalled {
		t.Error("Expected custom function to be called")
	}
	
	if _, ok := ds.(*mockOrderStatusDataSourceCustom); !ok {
		t.Error("Expected custom data source type")
	}
}

type mockOrderDataSourceCustom struct{}

func (m *mockOrderDataSourceCustom) Create(order daos.OrderDAO) error {
	return nil
}

func (m *mockOrderDataSourceCustom) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	return nil, nil
}

func (m *mockOrderDataSourceCustom) FindByID(id string) (daos.OrderDAO, error) {
	return daos.OrderDAO{}, nil
}

func (m *mockOrderDataSourceCustom) Update(order daos.OrderDAO) error {
	return nil
}

func (m *mockOrderDataSourceCustom) Delete(id string) error {
	return nil
}

type mockOrderStatusDataSourceCustom struct{}

func (m *mockOrderStatusDataSourceCustom) FindAll() ([]daos.OrderStatusDAO, error) {
	return nil, nil
}

func (m *mockOrderStatusDataSourceCustom) FindByID(id string) (daos.OrderStatusDAO, error) {
	return daos.OrderStatusDAO{}, nil
}