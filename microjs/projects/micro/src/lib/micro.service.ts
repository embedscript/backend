import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class MicroService {

  constructor() { }

  alert() {
    alert("hi")
  }
}
