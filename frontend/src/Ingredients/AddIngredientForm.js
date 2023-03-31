import React, { useState } from 'react';
import { Button, Modal, Form, Container } from 'react-bootstrap';

export function IngredientForm(props) {
    const [show, setShow] = useState(false);
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [image, setImage] = useState(null);
  
    const handleClose = () => setShow(false);
    const handleShow = () => setShow(true);
  
    const handleNameChange = (event) => {
      setName(event.target.value);
    }
  
    const handleDescriptionChange = (event) => {
      setDescription(event.target.value);
    }
  
    const handleImageChange = (event) => {
      setImage(event.target.files[0]);
    }
  
    const handleSubmit = (event) => {
      event.preventDefault();
      
      // Create a new FormData object and add the form data to it
      const formData = new FormData();
      formData.append('name', name);
      formData.append('description', description);
      formData.append('image', image);
  
      // Send a REST API request to your service
      fetch('https://your-service.com/api/ingredients', {
        method: 'POST',
        body: formData
      })
      .then(response => response.json())
      .then(data => {
        console.log(data);
        setShow(false);
      })
      .catch(error => {
        console.error(error);
      });
    }
  
    return (
      <Container>
        <Button variant="primary"  onClick={handleShow}>
          Add Ingredient
        </Button>
  
        <Modal show={show} onHide={handleClose}>
          <Modal.Header closeButton>
            <Modal.Title>Add Ingredient</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <Form onSubmit={handleSubmit}>
              <Form.Group controlId="formIngredientName">
                <Form.Label>Name</Form.Label>
                <Form.Control type="text" placeholder="Enter ingredient name" value={name} onChange={handleNameChange} />
              </Form.Group>
  
              <Form.Group controlId="formIngredientDescription">
                <Form.Label>Description</Form.Label>
                <Form.Control type="text" placeholder="Enter ingredient description" value={description} onChange={handleDescriptionChange} />
              </Form.Group>
  
              <Form.Group controlId="formIngredientImage">
                <Form.Label>Image</Form.Label>
                <Form.Control type="file" accept="image/*" onChange={handleImageChange} />
              </Form.Group>
            </Form>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={handleClose}>
              Cancel
            </Button>
            <Button variant="primary" onClick={handleSubmit}>
              Save
            </Button>
          </Modal.Footer>
        </Modal>
      </Container>
    );
  }