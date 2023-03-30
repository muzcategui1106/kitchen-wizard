import React, { useState, useEffect } from 'react';
import { Container, Row, Col, Image, Card } from 'react-bootstrap';


export function IngredientList() {
  const [ingredients, setIngredients] = useState([]);

  useEffect(() => {
    const headers = new Headers();
    headers.append('Cookie', document.cookie);
    console.log(document.cookie)

    fetch('http://localhost:8443/api/v1/ingredients?numItems=10', {
      credentials: 'include'
    })
      .then(response => response.json())
      .then(data => setIngredients(data));
  }, []);

    return (
      <Container>
        <Row>
          {ingredients.map((ingredient, index) => (
            <Col md={4} key={index} className="my-3">
              <Card>
                <Image src={ingredient.imageSrc} fluid />
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