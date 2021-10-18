import { Component } from '@angular/core';
import { LayoutService } from '../layout.service';
@Component({
  selector: 'jhi-sidebar',
  templateUrl: './sidebar.component.html',
  styleUrls: ['./sidebar.component.scss'],
})
export class SidebarComponent {
  constructor(private layoutService: LayoutService) {}

  toggleSidebar(): void {
    this.layoutService.toggleSidebar();
  }
}
