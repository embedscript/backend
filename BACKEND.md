# Backend

Backend is [Backend as a Service](https://www.cloudflare.com/learning/serverless/glossary/backend-as-a-service-baas/).

## Overview

Backend provides Backend as a Service for frontend developers. It's focused on giving devs the ability to build 
applications without having to worry about infrastructure. While M3O provides Micro as a Service for those 
who want to build Micro services, Backend goes a step further to focus on the frontend developers needs.

It's like Firebase https://firebase.google.com/products

## Features

- **Subscribe**: Subscription mailing lists
- **Database**: Document data storage
- **Files**: Basic file storage
- **Users**: User authentication and management
- **Notifications**: Push notifications via sms, email, websocket and mobile
- **Location**: Geolocation tracking
- **CRM**: Email marketing management
- **Functions**: Ability to define snippets of code to be executed
- **Hosting**: Hosting of static html/js apps

## Services

- Subscribe
- Database
- Document
- Push
- Users
- Location

## Document Storage

Document storage will be provided by a Document service. This is a layer of abstraction on top of the Store

  * github.com/micro/services/document
  * make use of github.com/micro/dev/model

## API

Backend provides a standard http/json API via `backend.m3o.dev`
