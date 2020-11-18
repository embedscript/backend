import { Component, OnInit } from '@angular/core';
import * as t from '../types';

const embeds: t.Embed[] = [
  {
    ID: 'blogging',
    Name: 'Blogging',
    Description:
      'Turn you static site into a dynamic blog with a few widgets. Lists posts, get a single post and edit them.',
    Available: true,
  },
  {
    ID: 'subscribe',
    Name: 'Subscribe',
    Description:
      'Put this widget on your site and start building your mailing list right now!',
    Available: false,
  },
  {
    ID: 'custom-form',
    Name: 'Custom form',
    Description:
      'Create custom forms that accept any kind of field or data and get them saved!',
    Available: false,
  },
  {
    ID: 'comments',
    Name: 'Comments',
    Description: 'Put a simple comments section on any content of your liking.',
    Available: false,
  },
  {
    ID: 'feeds',
    Name: 'Feeds',
    Description:
      'Embed your Wordpress or Medium blog posts on your static site. Edit on Wordpress, display on yours!',
    Available: false,
  },
  {
    ID: 'contact',
    Name: 'Contact',
    Description:
      'Get an email immediately when submits this form so you never miss out on leads.',
    Available: false,
  },
];

@Component({
  selector: 'app-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.css'],
})
export class ListComponent implements OnInit {
  embeds = embeds;
  constructor() {}

  ngOnInit(): void {

  }
}
