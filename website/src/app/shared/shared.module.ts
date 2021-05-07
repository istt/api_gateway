import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ToastrModule } from 'ngx-toastr';
import { ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
// + app specific
import { FormlyModule } from '@ngx-formly/core';
import { FormlyBootstrapModule } from '@ngx-formly/bootstrap';
import { FileValueAccessorDirective } from './util/file-value-accessor';
import { FormlyFieldFile } from './fields/file-type.component';
import { FormlyFileUploadComponent } from './fields/file-upload.type';
import { FormlyTabsetTypeComponent } from './fields/tabset.type';
import { RepeatTypeComponent } from './fields/repeat-section.type';
import { ButtonTypeComponent } from './fields/button.type';
import { AppValidators } from './util/app-validators';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { FaIconLibrary } from '@fortawesome/angular-fontawesome';
import { fas } from '@fortawesome/free-solid-svg-icons';

@NgModule({
  declarations: [
    FileValueAccessorDirective,
    FormlyFileUploadComponent,
    FormlyTabsetTypeComponent,
    RepeatTypeComponent,
    ButtonTypeComponent,
    FormlyFieldFile
  ],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    HttpClientModule,
    FontAwesomeModule,
    NgbModule,
    ToastrModule,
    FormlyBootstrapModule,
    FormlyModule.forRoot({
      validators: [{ name: 'ip', validation: AppValidators.IpValidator }],
      validationMessages: [
        { name: 'minlength', message: AppValidators.minlengthMessage },
        { name: 'maxlength', message: AppValidators.maxlengthMessage },
        { name: 'min', message: AppValidators.minMessage },
        { name: 'max', message: AppValidators.maxMessage },
        { name: 'minbytes', message: AppValidators.minbytesMessage },
        { name: 'maxbytes', message: AppValidators.maxbytesMessage },
        { name: 'pattern', message: AppValidators.patternMessage },
        { name: 'number', message: 'This field should be a number.' },
        { name: 'email', message: 'This field should be a valid email address.' },
        { name: 'datetimelocal', message: 'This field should be a date and time.' },
        { name: 'patternLogin', message: 'This field can only contain letters, digits and e-mail addresses.' },
        // + custom validators
        { name: 'ip', message: AppValidators.ipMessage },
        // - custom validators
        { name: 'required', message: 'This field is required.' }
      ],
      types: [
        { name: 'tabset', component: FormlyTabsetTypeComponent },
        { name: 'button', component: ButtonTypeComponent },
        { name: 'repeat', component: RepeatTypeComponent },
        { name: 'file-upload', component: FormlyFileUploadComponent },
      ],
    }),
  ],
  exports: [
    CommonModule,
    ReactiveFormsModule,
    HttpClientModule,
    FontAwesomeModule,
    NgbModule,
    ToastrModule,
    FormlyBootstrapModule,
    FormlyModule
  ]
})
export class SharedModule {
constructor(iconLibrary: FaIconLibrary) {
  iconLibrary.addIconPacks(fas);
}
}
