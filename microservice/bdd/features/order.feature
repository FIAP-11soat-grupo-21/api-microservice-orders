Feature: Order Management

  Scenario: Create a new order
    Given the order has valid data with customer ID and items
    When i send a request to create a new order
    Then the order should be created successfully