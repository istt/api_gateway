import { Component } from '@angular/core';
import { FieldArrayType } from '@ngx-formly/core';

@Component({
  selector: 'jhi-formly-tabset',
  template: `
    <div [ngClass]="field.className">
      <ul ngbNav #nav="ngbNav" [ngClass]="'nav-' + (to.type || 'tabs') + ' nav-' + (to.justify || 'justified')">
        <li ngbNavItem *ngFor="let panel of field.fieldGroup; let i = index" [disabled]="panel.templateOptions.disabled">
          <a ngbNavLink  [innerHtml]="panel.templateOptions.label"></a>
          <ng-template ngbNavContent>
            <div [ngClass]="panel.fieldGroupClassName || 'card'">
              <div [ngClass]="panel.className || 'card-body'">
                <formly-field class="col" [field]="f" *ngFor="let f of panel.fieldGroup"></formly-field>
              </div>
            </div>
          </ng-template>
        </li>
      </ul>
      <div [ngbNavOutlet]="nav"></div>
    </div>
  `
})
export class FormlyTabsetTypeComponent extends FieldArrayType {}
