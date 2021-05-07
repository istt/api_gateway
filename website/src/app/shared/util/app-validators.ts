import { AbstractControl } from '@angular/forms';

export class AppValidators {
  static jsonValidator(control: AbstractControl): any {
    try {
      JSON.parse(control.value);
    } catch (e) {
      return { invalidJson: e.message };
    }
    return null;
  }

  // Example validator from ngx-formly
  static IpValidator(control: AbstractControl): any {
    return control.value == null || /(\d{1,3}\.){3}\d{1,3}/.test(control.value) ? null : { ip: true };
  }

  static minlengthMessage(err: any): string {
    return `This field is required to be at least ${err.requiredLength} characters.`;
  }
  static maxlengthMessage(err: any): string {
    return `This field cannot be longer than ${err.requiredLength} characters.`;
  }
  static minMessage(err: any): string {
    return `This field should be at least ${err.min}.`;
  }
  static maxMessage(err: any): string {
    return `This field cannot be more than ${err.max}.`;
  }
  static minbytesMessage(err: any): string {
    return `This field should be at least ${err.minbytes} bytes.`;
  }
  static maxbytesMessage(err: any): string {
    return `This field cannot be more than ${err.maxbytes} bytes.`;
  }
  static patternMessage(err: any, field: any): string {
    return `${field.templateOptions.label} is should follow pattern for ${err.requiredPattern}.`;
  }
  static ipMessage(err: any, field: any): string {
    return `"${field.formControl.value}" is not a valid IP Address`;
  }
}
