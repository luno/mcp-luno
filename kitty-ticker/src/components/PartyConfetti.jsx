import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import PropTypes from 'prop-types';

const ConfettiPiece = styled.div`
  position: fixed;
  width: ${props => props.size}px;
  height: ${props => props.size}px;
  background-color: ${props => props.color};
  border-radius: ${props => props.shape === 'circle' ? '50%' : '0'};
  top: -20px;
  left: ${props => props.left}%;
  animation: fall ${props => props.duration}s linear infinite;
  animation-delay: ${props => props.delay}s;
  z-index: -1;
  transform: rotate(${props => props.rotation}deg);
  opacity: 0.7;

  @keyframes fall {
    0% {
      transform: translateY(-20px) rotate(0deg);
    }
    100% {
      transform: translateY(100vh) rotate(720deg);
    }
  }
`;

const colors = [
  '#ff69b4', // Hot Pink
  '#00ffff', // Cyan
  '#ff6347', // Tomato
  '#7fffd4', // Aquamarine
  '#ff8c00', // Dark Orange
  '#9370db', // Medium Purple
  '#32cd32', // Lime Green
  '#ffd700'  // Gold
];

const shapes = ['square', 'circle'];

const PartyConfetti = ({ count = 20 }) => {
  const [confetti, setConfetti] = useState([]);

  useEffect(() => {
    const pieces = [];

    for (let i = 0; i < count; i++) {
      pieces.push({
        id: i,
        color: colors[Math.floor(Math.random() * colors.length)],
        size: Math.random() * 10 + 5,
        left: Math.random() * 100,
        duration: Math.random() * 5 + 3,
        delay: Math.random() * 5,
        rotation: Math.random() * 360,
        shape: shapes[Math.floor(Math.random() * shapes.length)]
      });
    }

    setConfetti(pieces);
  }, [count]);

  return (
    <>
      {confetti.map(piece => (
        <ConfettiPiece
          key={piece.id}
          color={piece.color}
          size={piece.size}
          left={piece.left}
          duration={piece.duration}
          delay={piece.delay}
          rotation={piece.rotation}
          shape={piece.shape}
        />
      ))}
    </>
  );
};

PartyConfetti.propTypes = {
  count: PropTypes.number
};

export default PartyConfetti;
