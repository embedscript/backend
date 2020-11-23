import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MicroComponent } from './micro.component';

describe('MicroComponent', () => {
  let component: MicroComponent;
  let fixture: ComponentFixture<MicroComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MicroComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MicroComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
