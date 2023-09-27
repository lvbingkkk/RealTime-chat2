import { Component, OnDestroy, OnInit } from '@angular/core';
import { SocketService } from "./socket.service";


@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent  implements OnInit, OnDestroy {


  public messages: Array<any>;
  public chatBox: string;

  public constructor(private socket: SocketService) {
      this.messages = [];
      this.chatBox = "";
  }

  public ngOnInit() {
      this.socket.getEventListener().subscribe(event => {
        console.log('nginit event:',event)

        // console.log('nginit: event.type',event.type)

          if(event.type == "message") {
            console.log('nginit: event.data',event.data)

              let data = event.data.content;
              if(event.data.sender) {
                  data = event.data.sender + ": " + data;
              }
              this.messages.push(data);
          }
          if(event.type == "close") {
            // this.messages.push("/The socket connection has been closed");
            this.messages.push("/您的时空隧道连接已关闭!!");
        }
          if(event.type == "open") {
            // this.messages.push("/The socket connection has been established");
            this.messages.push("/您的时空隧道连接已建立!!");
        }
      });
  }

  public ngOnDestroy() {
      this.socket.close();
  }

  public send() {
      if(this.chatBox) {
          this.socket.send(this.chatBox);
          this.chatBox = "";
      }
  }

  public isSystemMessage(message: string) {
      return message.startsWith("/") ? "<strong>" + message.substring(1) + "</strong>" : message;
  }
}
