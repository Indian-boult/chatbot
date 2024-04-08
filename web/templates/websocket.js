const ws = new WebSocket('ws://localhost:8080/ws');
const chatBox = document.getElementById('chat-box');
const chatInput = document.getElementById('chat-input');

ws.onmessage = function(event) {
    // Check if the message contains an image URL
    if (event.data.startsWith("Here is your image URL: ")) {
        const imageUrl = event.data.replace("Here is your image URL: ", "");

        const fullImageUrl = `/uploads/${imageUrl.split('uploads/')[1]}`;
        const imageElement = document.createElement('img');
        imageElement.src = fullImageUrl;
        imageElement.style.maxWidth = '200px'; // Set image size as needed
        imageElement.style.borderRadius = '10px'; // Optional styling
        chatBox.appendChild(imageElement);
    } else {
        const message = document.createElement('div');
        message.textContent = event.data;
        chatBox.appendChild(message);
    }
};

function sendMessage() {
    const message = chatInput.value;
    ws.send(message);
    chatInput.value = ''; // Clear input after sending
}
