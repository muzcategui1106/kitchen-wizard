import React, { useState } from 'react';
import { Button, Modal, Form, Container, Alert  } from 'react-bootstrap';
import axios from 'axios';

export function IngredientForm(props) {
  const [show, setShow] = useState(false);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [image, setImage] = useState(null);
  const [apiErrorMessage, setApiErrorMessage] = useState('');

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

    // Read the file as a binary buffer
    const fileReader = new FileReader();
    fileReader.readAsArrayBuffer(image);

    fileReader.onload = async () => {
      const imageBytes = new Uint8Array(fileReader.result);


      console.log(imageBytes)
      const data = {
        ingredient: {
          name: name,
          description: description
        },
        image: Array.from(imageBytes)

      };

      const options = {
        url: 'http://localhost:8443/api/v1/ingredient',
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Content-Length': data.length
        },
        data: JSON.stringify(data),
        withCredentials: true
      };

      axios(options)
        .then((response) => {
          if (response.status == 200 ) {
            setShow(false); // close the modal form
          } else {
            console.log(`statusCode: ${response.status}`);
            console.log(response.data);
            setApiErrorMessage(response.data.error)
          }
        })
        .catch((error) => {
          setApiErrorMessage(JSON.stringify(error.response.data.error));
          console.error(error);
        });
    };


  }

  return (
    <Container>
      <Button variant="primary" onClick={handleShow}>
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
          {apiErrorMessage && <Alert variant="danger" className="mt-3">{apiErrorMessage}</Alert>}
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