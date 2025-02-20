const WebSocket = require('ws');

const socket = new WebSocket('ws://localhost:6969/ws'); // Change URL if needed

socket.onopen = () => {
    console.log('Connected to WebSocket server');
};

socket.onmessage = (event) => {
    try {
        const data = JSON.parse(event.data);
        console.log('Received event:', data);
    } catch (err) {
        console.error('Error parsing WebSocket message:', err);
    }
};

socket.onclose = () => {
    console.log('WebSocket connection closed');
};

socket.onerror = (error) => {
    console.error('WebSocket error:', error);
};
