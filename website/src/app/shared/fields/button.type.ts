import { Component } from '@angular/core';
import { FieldType } from '@ngx-formly/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { createRequestOption, plainToFlattenObject } from '../util/request-util';
import * as _ from 'lodash';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';

@Component({
  selector: 'jhi-formly-button',
  template: `
    <div class="my-auto">
      <button [type]="to.type || 'button'" [ngClass]="'btn btn-' + (to.btnType || 'outline-primary')" (click)="onClick($event)">
        {{ to.label }}
      </button>
    </div>
  `
})
export class ButtonTypeComponent extends FieldType {
  constructor(protected httpClient: HttpClient) {
    super();
  }
  onClick($event: any): void {
    if (this.to.onClick) {
      this.to.onClick($event);
    } else if (this.to.apiEndpoint) {
      this.createRequest()
        .pipe(
          filter(res => res.ok),
          map(res => res.body)
        )
        .subscribe(
          res => this.processResponse(res),
          err => console.error(this.to.errorMsg || err.message)
        );
    }
  }

  createRequest(): Observable<HttpResponse<any>> {
    const params = createRequestOption(_.omitBy(plainToFlattenObject(this.to.params), _.isNull));
    const body = _.omitBy(this.to.body, _.isNull);
    if (this.to.method) {
      return this.httpClient.request<HttpResponse<any>>(this.to.method,  this.to.apiEndpoint, {
        params,
        body,
        observe: 'response'
      });
    } else if (_.isEmpty(body)) {
      return this.httpClient.get<HttpResponse<any>>( this.to.apiEndpoint, {
        params,
        observe: 'response'
      });
    } else {
      return this.httpClient.post<HttpResponse<any>>( this.to.apiEndpoint, body, {
        params,
        observe: 'response'
      });
    }
  }

  processResponse(res: any): void {
    this.formControl.setValue(res);
  }
}
