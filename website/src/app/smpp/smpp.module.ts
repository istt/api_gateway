import { NgModule } from '@angular/core';
import { SharedModule } from 'app/shared/shared.module';

import { RouterModule } from '@angular/router';
import { routes } from './smpp.routing';

import { EsmeComponent } from './esme/esme.component';

@NgModule({
  declarations: [
    EsmeComponent
  ],
  imports: [
    SharedModule,
    RouterModule.forChild(routes)
  ]
})
export class SmppModule { }
