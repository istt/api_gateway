import { Routes } from '@angular/router';

// app components
import { EsmeComponent } from './esme/esme.component';

export const routes: Routes = [
  {
    path: '',
    redirectTo: 'esme-profiles',
    pathMatch: 'full',
  },
  {
    path: 'esme-profiles',
    component: EsmeComponent
  }
];
