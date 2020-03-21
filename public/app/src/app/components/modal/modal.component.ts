import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-modal',
  templateUrl: './modal.component.html',
  styleUrls: ['./modal.component.scss']
})
export class ModalComponent implements OnInit {
  visible = false
  constructor() { }

  ngOnInit(): void {
  }

  public show() {
    this.visible = true
  }
  public hide() {
    this.visible = false
  }

}
