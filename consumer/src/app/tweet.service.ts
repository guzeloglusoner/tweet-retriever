import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';
import { map } from 'rxjs/operators';
import { WebSocketService } from './websocket.service';

export interface Message {
    data: string;
}

@Injectable({
    providedIn: 'root'
})
export class TweetService {
    public messages: Subject<Message>;

    constructor(_websocketService: WebSocketService) {
        this.messages = <Subject<Message>>_websocketService.connect('ws://localhost:9090/ws').pipe(map(
            (response: MessageEvent): Message => {
                /*const responseJSON = JSON.parse(response.data);*/
                return {
                    data: response.data,
                };
            }
        ));
     }
}
