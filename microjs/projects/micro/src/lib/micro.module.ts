import { NgModule } from '@angular/core';
import { MicroComponent } from './micro.component';
import { LoginComponent } from './login/login.component';

@NgModule({
  declarations: [MicroComponent, LoginComponent],
  imports: [LoginComponent],
  exports: [MicroComponent],
})
export class MicroModule {}
