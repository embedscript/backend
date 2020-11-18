import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { ListComponent } from './list/list.component';
import { SingleComponent } from './single/single.component';

const routes: Routes = [
  {
    path: ':id',
    component: SingleComponent,
  },
  {
    path: '',
    component: ListComponent,
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
