import React, { useState, useEffect } from "react"; 
import { useHistory } from "react-router-dom"
import Nav from "react-bootstrap/Nav";

export function NavBarLoginButton() {

    const history = useHistory();
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    var loginWindow;
    
    useEffect(() => {
      window.addEventListener('message', handleMessage);
    
      return () => {
        window.removeEventListener('message', handleMessage);
      };
    }, []);
    
    function handleMessage(event) {
      console.log('Received message:', event);
      console.log(document.cookie);
        
      console.log(event)
      if (event.origin !== 'http://localhost:3000') {
        console.log('Invalid origin:', event.origin);
        return;
      }
    
      const token = event.data.token;
    
      if (token) {
        document.cookie = `token=${token}; path=/;`;
        setIsLoggedIn(true);
        history.push('/');
        loginWindow.close()
      } else {
        console.log(event.data)
      }
    }
    
    function handleLogin() {
      console.log('before event');
      console.log(document.cookie);
      loginWindow = window.open('http://localhost:8443/auth/v1/login', '_blank');
    }
    
    function handleLogout() {
      // Make an API call to logout
      // Set loggedIn to false
      // Redirect the user to the login page
      history.push("/login");
    }
  
    return (
      <Nav className="ml-auto">
        {isLoggedIn ? (
          <Nav.Link onClick={handleLogout}>Logout</Nav.Link>
        ) : (
          <Nav.Link onClick={handleLogin}>Login</Nav.Link>
        )}
      </Nav>
    );
  }