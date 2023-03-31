import React, { useState, useEffect } from "react";
import { useHistory } from "react-router-dom"
import Nav from "react-bootstrap/Nav";

export function NavBarLoginButton() {

  const history = useHistory();
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  var loginWindow;
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    // Check if the 'kitchenwizard' cookie is present
    if (document.cookie.match(/(^|;)\s*kitchenwizard=/)) {
      setIsLoggedIn(true)
    } else {
      setIsLoggedIn(false)
    }
  }, []);

  function handleLogout() {
    // Make an API call to logout
    // Set loggedIn to false
    // Redirect the user to the login page
    history.push("/login");
    // Clear the 'kitchenwizard' cookie by setting its value to an empty string
    // and its expiration date to a date in the past
    document.cookie = 'kitchenwizard=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
  }

  const handleClick = () => {
    setIsLoading(true);
    fetch('http://localhost:8443/auth/v1/login')
      .then(response => {
        if (response.ok) {
          const newTab = window.open(response.url, '_blank');
          newTab.focus();
          setIsLoading(false);
          newTab.addEventListener('load', event => {
            console.log("got event")
            console.log(event.data)
          });
        } else {
          console.error(`API returned ${response.status} ${response.statusText}`);
          setIsLoading(false);
        }
      })
      .catch(error => console.error(error));
  };

  return (
    <Nav className="ml-auto">
      {isLoggedIn ? (
        <Nav.Link onClick={handleLogout}>Logout</Nav.Link>
      ) : (
        <Nav.Link onClick={handleClick}>Login</Nav.Link>
      )}
    </Nav>
  );
}