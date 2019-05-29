import { Component, OnInit } from '@angular/core';
import { TweetService, Message } from './tweet.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  title = 'consumer';
  messages: Message[] = [];

  constructor(tweetService: TweetService) {
    tweetService.messages.subscribe(msg => {
      this.messages.push(msg);
    });
   }

  ngOnInit() {
  }
}
