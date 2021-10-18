import { Component, OnInit, RendererFactory2, Renderer2 } from '@angular/core';
import { SwUpdate } from '@angular/service-worker';
import { Title } from '@angular/platform-browser';
import { Router, ActivatedRouteSnapshot, NavigationEnd, NavigationError } from '@angular/router';
import { TranslateService, LangChangeEvent } from '@ngx-translate/core';
import * as dayjs from 'dayjs';
import { LayoutService } from '../layout.service';
import { AccountService } from 'app/core/auth/account.service';

@Component({
  selector: 'jhi-main',
  templateUrl: './main.component.html',
})
export class MainComponent implements OnInit {
  private renderer: Renderer2;

  constructor(
    private swUpdate: SwUpdate,
    private layoutService: LayoutService,
    private accountService: AccountService,
    private titleService: Title,
    private router: Router,
    private translateService: TranslateService,
    rootRenderer: RendererFactory2
  ) {
    this.renderer = rootRenderer.createRenderer(document.querySelector('html'), null);
    swUpdate.available.subscribe(() => {
      if (confirm(`New version available. Load New Version?`)) {
        swUpdate.activateUpdate().then(() => document.location.reload());
      }
    });
  }

  ngOnInit(): void {
    // try to log in automatically
    this.accountService.identity().subscribe();

    this.router.events.subscribe(event => {
      if (event instanceof NavigationEnd) {
        this.updateTitle();
      }
      if (event instanceof NavigationError && event.error.status === 404) {
        this.router.navigate(['/404']);
      }
    });

    this.translateService.onLangChange.subscribe((langChangeEvent: LangChangeEvent) => {
      this.updateTitle();
      dayjs.locale(langChangeEvent.lang);
      this.renderer.setAttribute(document.querySelector('html'), 'lang', langChangeEvent.lang);
    });
  }

  getClasses(): any {
    const classes = {
      'pinned-sidebar': this.layoutService.getSidebarStat().isSidebarPinned,
      'toggeled-sidebar': this.layoutService.getSidebarStat().isSidebarToggeled,
    };
    return classes;
  }

  toggleSidebar(): void {
    this.layoutService.toggleSidebar();
  }

  isAuthenticated(): boolean {
    return this.accountService.isAuthenticated();
  }

  private getPageTitle(routeSnapshot: ActivatedRouteSnapshot): string {
    let title: string = routeSnapshot.data['pageTitle'] ?? '';
    if (routeSnapshot.firstChild) {
      title = this.getPageTitle(routeSnapshot.firstChild) || title;
    }
    return title;
  }

  private updateTitle(): void {
    let pageTitle = this.getPageTitle(this.router.routerState.snapshot.root);
    if (!pageTitle) {
      pageTitle = 'global.title';
    }
    this.translateService.get(pageTitle).subscribe(title => this.titleService.setTitle(title));
  }
}
