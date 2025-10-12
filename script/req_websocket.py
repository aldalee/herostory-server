import websocket

URL = "ws://localhost:12345/websocket"


def send_single_websocket_message():
    try:
        ws = websocket.create_connection(URL)
        print("Connected to WebSocket Server")

        msg = "Hello World"
        ws.send(msg)
        print(f"Sent: {msg}")

        # response = ws.recv()
        # print(f"Received: {response}")

        ws.close()
        print("WebSocket connection closed")
    except Exception as e:
        print(f"WebSocket error occurred: {e}")


if __name__ == "__main__":
    send_single_websocket_message()
