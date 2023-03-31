import React from "react"; 
import Navbar from "react-bootstrap/Navbar";
import Nav from "react-bootstrap/Nav";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import NavDropdown from "react-bootstrap/NavDropdown";
import { NavBarLoginButton } from '../Login/Login';
import { IngredientListPanel } from '../Ingredients/IngredientList.js'


export function TopNavigationBar() {
  return (
    <Router>
      <Navbar bg="dark" variant="dark" expand="lg">
          <Navbar.Brand href="#home">Kitchen Wizard</Navbar.Brand>
          <Navbar.Toggle aria-controls="basic-navbar-nav" />
          <Navbar.Collapse id="basic-navbar-nav">
            <Nav className="me-auto">
              <Nav.Link href="#Home">Home</Nav.Link>
              <Nav.Link href="/ingredients">Ingredients</Nav.Link>
              <NavDropdown title="Dropdown" id="basic-nav-dropdown">
                <NavDropdown.Item href="#action/3.1">Action</NavDropdown.Item>
                <NavDropdown.Item href="#action/3.2">
                  Another action
                </NavDropdown.Item>
                <NavDropdown.Item href="#action/3.3">Something</NavDropdown.Item>
                <NavDropdown.Divider />
                <NavDropdown.Item href="#action/3.4">
                  Separated link
                </NavDropdown.Item>
              </NavDropdown>
            </Nav>
            <NavBarLoginButton />
          </Navbar.Collapse>
      </Navbar>
      <Switch>
        <Route path="/ingredients">
          <IngredientListPanel />
        </Route>
      </Switch>
    </Router>

    
  );
}