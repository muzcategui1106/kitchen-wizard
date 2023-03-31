import React, { useState, useEffect } from 'react';
import { Container, Row, Col, Image, Card } from 'react-bootstrap';
import { IngredientForm } from './AddIngredientForm'

export function IngredientListPanel() {
  return (
    <Container>
      <IngredientForm />
      <IngredientList />
    </Container>
  )
}

function IngredientList() {
  const [ingredients, setIngredients] = useState([]);


  async function downloadImage(imageUrl) {
    try {
      const response = await fetch(imageUrl);
      const blob = await response.blob();
      const url = URL.createObjectURL(blob);
      return url;
    } catch (error) {
      console.error('Error downloading image:', error);
      throw error;
    }
  }

  useEffect(() => {
    const headers = new Headers();
    headers.append('Cookie', document.cookie);
    console.log(document.cookie)

    fetch('http://localhost:8443/api/v1/ingredients?numItems=10', {
      credentials: 'include'
    })
      .then(response => response.json())
      .then(data => {
        // For each ingredient, download its image and set its URL in state
        const promises = data.map(async (ingredient) => {
          const imageUrl = `http://kitchen-wizard.store.s3.local.uzcatm-skylab.com/kitchen-wizard/ingredients/${ingredient.name}.jpg`;
          const url = await downloadImage(imageUrl);
          return {
            ...ingredient,
            imageUrl: url
          };
        });

        Promise.all(promises).then((data) => setIngredients(data));
      });
    }, []);

  return (
    <Container>
      <Row>
        {ingredients.map((ingredient, index) => (
          <Col md={4} key={index} className="my-3">
            <Card>
              <Image
                src={ingredient.imageUrl}
                alt={ingredient.name}
                fluid
              />
              <Card.Body>
                <Card.Title>{ingredient.name}</Card.Title>
                <Card.Text>{ingredient.description}</Card.Text>
              </Card.Body>
            </Card>
          </Col>
        ))}
      </Row>
    </Container>
  );
};