import { HttpParams } from '@angular/common/http';
import * as _ from 'lodash';

export const createRequestOption = (req?: any): HttpParams => {
  let options: HttpParams = new HttpParams();
  if (req) {
    _.each(req, (val, key) => {
      if (key !== 'sort') {
        if (_.isArray(val)) {
          _.each(val, v => (options = options.append(key, v)));
        } else {
          options = options.set(key, req[key]);
        }
      }
    });
    if (req.sort) {
      req.sort.forEach(val => {
        options = options.append('sort', val);
      });
    }
  }
  return options;
};

export const plainToFlattenObject = (object: any) => {
  const result = {};

  function flatten(obj, prefix = '') {
    _.forEach(obj, (value, key) => {
      if (_.isObject(value)) {
        flatten(value, `${prefix}${key}.`);
      } else {
        result[`${prefix}${key}`] = value;
      }
    });
  }

  flatten(object);

  return result;
};
