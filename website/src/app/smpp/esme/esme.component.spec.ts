import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EsmeComponent } from './esme.component';

describe('EsmeComponent', () => {
  let component: EsmeComponent;
  let fixture: ComponentFixture<EsmeComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ EsmeComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(EsmeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
