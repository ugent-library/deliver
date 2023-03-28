import BSN from 'bootstrap.native/dist/bootstrap-native-v4'
import bsnPopper from '../lib/bsn_popper.js'
import { Controller } from "@hotwired/stimulus";

export default class extends Controller {
    connect() {
        BSN.initCallback(this.element)
        bsnPopper(this.element)
    }
}