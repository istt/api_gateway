import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { HttpErrorResponse, HttpHeaders, HttpResponse, HttpClient } from '@angular/common/http';
import { ActivatedRoute, Router, NavigationEnd } from '@angular/router';
import { FormGroup } from '@angular/forms';
import { Subscription, forkJoin, of, empty, from  } from 'rxjs';
import { filter, map, catchError, tap, concatAll } from 'rxjs/operators';
import { createRequestOption } from 'app/shared/util/request-util';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
// + search
import * as _ from 'lodash';
import * as jsyaml from 'js-yaml';
// + websocket support
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';

@Component({
  selector: 'app-esme',
  templateUrl: './esme.component.html',
  styleUrls: ['./esme.component.css']
})
export class EsmeComponent implements OnInit, OnDestroy {
  // + profiles
  profiles: any[];
  // + config
  form = new FormGroup({});
  fields: any;
  settings: any = {
    tlvList: []
  };
  _ = _;
  messages: any[] = [];
  smsSet: any[] = [];

  error: any;
  success: any;
  statusMsg: any;
  isWs = false;
  wsEndpoint: string;
  wsTopic: string = 'smpp-simulator';
  wsSocket: WebSocketSubject<any>;

  // + delete Modal
  @ViewChild('detailModal', { static: true }) detailModal: any;
  modalModel: any = {};

  constructor(
    protected httpClient: HttpClient,
    private modalService: NgbModal
  ) {
    this.settings = {};
  }

  ngOnInit() {
    this.loadSettings();
    this.loadForm();
    this.loadProfiles();
    this.openWs();
  }
  ngOnDestroy() {
    if (this.wsSocket) {
      this.wsSocket.unsubscribe();
    }
  }

  openWs(): void {
    this.isWs = true;
    this.wsEndpoint = 'ws://' + window.location.host + '/ws';
    this.wsSocket = webSocket(this.wsEndpoint);
    this.wsSocket.subscribe(
      activity => this.parseMsg(activity), // Called whenever there is a message from the server.
      err => console.error('Cannot parseMsg ' + this.wsTopic, err), // Called if at any point WebSocket API signals some kind of error.
      () => this.isWs = false // Called when connection is closed (for whatever reason).
    );
  }
  parseMsg(activity: any): void {
    if (!activity) return;
    if (activity.message) {
      try {
        const data = JSON.parse(activity.message);
        if (data.Cookie) {
          this.parseMsg(data.Cookie)
        }
      } catch (e) {
        console.error(e);
      }
    } else if (activity.codeStatus) {
      if (activity.codeStatus === 2) {
        this.statusMsg = _.get(activity, 'data.info');
      } else if (activity.data) {
        this.messages.unshift(activity.data);
      }
    }

  }

  protected onError(errorMessage: string) {
    console.error(errorMessage, null, null);
  }

  // Load current settings from backend
  loadSettings() {
    this.httpClient.get('api/settings', { observe: 'response' })
      .pipe(
        tap(res => console.log(res)),
        filter(res => res.ok),
        map(res => res.body)
      ).subscribe(res => {
        this.settings = res ? res : {};
        if (!this.settings.tlvList) {
          _.set(this.settings, 'tlvList', []);
        }
      }, err => this.onError(err.message));
  }
  // Display form metadata from YAML
  loadForm() {
    this.httpClient.get('assets/forms/smpp-simulator.yaml', { observe: 'response', responseType: 'text' })
      .pipe(
        filter(res => res.ok),
        map(res => jsyaml.load(res.body))
      ).subscribe(res => this.fields = res, err => this.onError(err.message));
  }
  // Load profiles from consul
  loadProfiles() {
    this.httpClient.get('api/esme-profiles', { params: createRequestOption({ key: 'smpp-simulator' }), observe: 'response', responseType: 'text' })
      .pipe(
        filter(res => res.ok),
        map(res => jsyaml.load(res.body)),
        catchError(err => of([]))
      ).subscribe(res => this.profiles = res);
  }
  // Save profiles to backend
  saveProfiles() {
    this.httpClient.post('api/esme-profiles', jsyaml.dump(this.profiles), { params: createRequestOption({ key: 'smpp-simulator' }), observe: 'response', responseType: 'text' })
      .pipe(
        filter(res => res.ok),
        map(res => jsyaml.load(res.body))
      ).subscribe(res => this.profiles = res, err => this.onError(err.message));
  }
  // Select one of the profiles to load into current setting
  selectProfile(profile) {
    this.settings = profile;
  }
  addProfile() {
    this.profiles.push(_.assign({}, this.settings));
    this.saveProfiles();
  }
  deleteProfile(profile) {
    _.pull(this.profiles, profile);
    this.saveProfiles();
  }

  save() {
    this.httpClient
      .post('api/settings', this.settings, {
        observe: 'response'
      })
      .pipe
      // tap(res => this.jhiAlertService.success('successfully save settings'))
      ()
      .subscribe(res => (this.settings = res.body), err => this.onError(err.message));
  }

  startASession() {
    this.httpClient
      .post('api/start-session', this.settings, { observe: 'response' })
      .pipe(
        // tap(res => this.jhiAlertService.success('startASession')),
        map(res => res.body)
      )
      .subscribe((res: any) => !this.isWs && this.parseMsg(res), err => this.onError(err.message));
  }

  stopASession() {
    this.httpClient
      .get('api/stop-session', { observe: 'response' })
      .pipe(
        // tap(res => this.jhiAlertService.success('stopASession')),
        map(res => res.body)
      )
      .subscribe((res: any) => !this.isWs && this.parseMsg(res), err => this.onError(err.message));
  }

  refreshState() {
    this.httpClient
      .get('api/refresh-state', { observe: 'response' })
      .pipe(
        // tap(res => this.jhiAlertService.success('refreshState')),
        map(res => res.body)
      )
      .subscribe((res: any) => !this.isWs && this.parseMsg(res), err => this.onError(err.message));
  }

  sendBadPacket() {
    this.httpClient
      .get('api/send-bad-packet', { observe: 'response' })
      .pipe(
        // tap(res => this.jhiAlertService.success('sendBadPacket')),
        map(res => res.body)
      )
      .subscribe((res: any) => !this.isWs && this.parseMsg(res), err => this.onError(err.message));
  }

  submitMessage() {
    this.httpClient
      .post('api/submit-message', this.settings, { observe: 'response' })
      .pipe(
        // tap(res => this.jhiAlertService.success('submitMessage')),
        map(res => res.body)
      )
      .subscribe((res: any) => !this.isWs && this.parseMsg(res), err => this.onError(err.message));
  }

  bulkSendingRandom() {
    this.httpClient
      .get('api/bulk-sending-random', { observe: 'response' })
      .pipe(
        // tap(res => this.jhiAlertService.success('bulkSendingRandom')),
        map(res => res.body)
      )
      .subscribe((res: any) => !this.isWs && this.parseMsg(res), err => this.onError(err.message));
  }

  stopBulkSending() {
    this.httpClient
      .get('api/stop-bulk-sending', { observe: 'response' })
      .pipe(
        // tap(res => this.jhiAlertService.success('stopBulkSending')),
        map(res => res.body)
      )
      .subscribe((res: any) => !this.isWs && this.parseMsg(res), err => this.onError(err.message));
  }
  refresh() {
    this.messages = [];
  }

  viewDetail(m: any) {
      this.modalModel = m;
      this.modalService.open(this.detailModal, { size: 'lg' }).result.then(
        () => this.modalService.dismissAll(),
        () => this.modalService.dismissAll()
      );
  }

  convertHex(): void {
    this.httpClient.post('api/convert-text', _.pick(this.settings, ['charset', 'messageText']), { observe: 'response', responseType: 'text'})
    .pipe(
      map(res => res.body)
    )
    .subscribe(shortMsgHex => this.form.patchValue({ shortMsgHex }));
  }

  addToSet(): void {
    this.smsSet.push(_.assign({}, this.settings));
  }

  deleteSms(idx: number): void {
    _.pullAt(this.smsSet, idx);
  }

  sendSmsSet(): void {
    from(this.smsSet)
      .pipe(
        map(sms => this.httpClient.post('api/submit-message', sms, { observe: 'response' })),
        concatAll()
      )
      .subscribe();
  }
}
