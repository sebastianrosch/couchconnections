import { Component, OnInit } from '@angular/core';
import { AuthService } from 'src/app/auth/auth.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {

  constructor(public auth: AuthService) {}

  ngOnInit() {
  }

  /**
   * Scroll button event handler
   */
  scrollToCalendar(): void {
    window.scrollBy({
      top: 600,
      left: 0,
      behavior: 'smooth'
    });
  }

}
