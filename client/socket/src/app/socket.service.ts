import { Injectable, EventEmitter} from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class SocketService {

  private socket: WebSocket;
  private listener: EventEmitter<any> = new EventEmitter();


  constructor() {
    console.log('new websocket')
    this.socket = new WebSocket("ws://localhost:12345/myws");
    //这是 ngrok 的穿透代理, 花生壳的tcp 代理就不行!!!why
    // cpolar 穿透也行 ./cpolar tcp 12345
    // this.socket = new WebSocket( "ws://6.tcp.ngrok.io:13475/myws");

    this.socket.onopen = event => {
        this.listener.emit({"type": "open", "data": event});
    }
    this.socket.onclose = event => {
        this.listener.emit({"type": "close", "data": event});
    }
    this.socket.onmessage = event => {
        this.listener.emit({"type": "message", "data": JSON.parse(event.data)});
    }
  }

  public send(data: string) {
    this.socket.send(data);
  }

  public close() {
    this.socket.close();
  }

  public getEventListener() {
    return this.listener;
  }

}
