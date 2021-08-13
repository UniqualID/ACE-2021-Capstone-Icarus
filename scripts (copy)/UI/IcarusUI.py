#!/usr/bin/env python
import gi
import subprocess
import sys

gi.require_version("Gtk", "3.0")
from gi.repository import Gtk


class IcarusUI:
    def __init__(self):
        self.gladefile = "IcarusUI.glade"
        self.builder = Gtk.Builder()
        self.builder.add_from_file(self.gladefile)
        self.builder.connect_signals(self)
        self.window = self.builder.get_object("window1")
        self.window.show()
        self.IP = self.builder.get_object("IP")
        self.ID = self.builder.get_object("ID")
        self.NAME = self.builder.get_object("TYPE")

    def on_window1_destroy(self, object, data=None):
        Gtk.main_quit()
        return

    def button_send(self, button, data=None):
        file = self.builder.get_object("ROUTES").get_text()
        subprocess.Popen(["./route", file, self.IP.get_text()])
        return

    def button_refuelall(self, button, data=None):
        subprocess.Popen(["./refuelall", self.IP.get_text()])
        return

    def button_follow(self, button, data=None):
        print("NOT ADDED")
        return

    def button_followfire(self, button, data=None):
        action = "1"
        payid = "0"

        if self.builder.get_object("LANCE").get_active():
            payid = "3"
        if self.builder.get_object("CAMERA").get_active():
            payid = "4"
        if self.builder.get_object("PHOSPHEX").get_active():
            payid = "7"
        if self.builder.get_object("PHOSPHEXREM").get_active():
            payid = "8"
        if self.builder.get_object("ANTIMATTER").get_active():
            payid = "10"
        if self.builder.get_object("SEEKER").get_active():
            payid = "15"
        if payid == "0":
            return

        tgt = self.builder.get_object("TGT_ID").get_text()
        if tgt == "":
            tgt = "0"

        id = self.builder.get_object("ID").get_text()
        if id == "":
            return
        id = id.zfill(3)
        veh = self.builder.get_object("TYPE").get_active_id() + "-" + id
        subprocess.Popen(
            ["./followFire", veh, payid, action, tgt,
             self.IP.get_text()])

    def button_emgland(self, button, data=None):
        id = self.builder.get_object("ID").get_text()
        if id == "":
            return
        id = id.zfill(3)
        veh = self.builder.get_object("TYPE").get_active_id() + "-" + id
        subprocess.Popen(["./emergencyland", veh, self.IP.get_text()])
        return

    def button_status_all(self, button, data=None):
        subprocess.Popen(["./status_all", self.IP.get_text()])
        return

    def button_connect(self, button, data=None):
        if not self.builder.get_object("ADMIN").get_active():
            return
        if self.builder.get_object("REAL").get_active():
            subprocess.Popen(
                ["./addallvehicles", "realassets.csv",
                 self.IP.get_text()])
        elif self.builder.get_object("TEST").get_active():
            subprocess.Popen(
                ["./addallvehicles", "testassets.csv",
                 self.IP.get_text()])
        else:
            return

    def button_disconnect(self, button, data=None):
        if not self.builder.get_object("ADMIN").get_active():
            return
        subprocess.Popen(["./removeallvehicles", self.IP.get_text()])
        return

    def button_goto(self, button, data=None):
        gps = self.builder.get_object("GPS").get_text()
        gps = gps.replace(",", " ")
        gps = gps.split()

        if gps[0] == "ELVENKING-AB":
            gps = ["50.6533", "-68.7170"]
        elif gps[0] == "EDORAS-AB":
            gps = ["49.1819", "-68.3639"]
        elif gps[0] == "ETHRING-AB":
            gps = ["49.6165", "-67.7075"]
        elif gps[0] == "EDHELLOND-AB":
            gps = ["50.2223", "-66.2656"]
        elif gps[0] == "SHELOB-AB":
            gps = ["50.5572", "-63.4119"]
        elif gps[0] == "SMAUG-AB":
            gps = ["50.8405", "-62.4290"]
        elif gps[0] == "SCATHA-AB":
            gps = ["50.2306", "-62.1124"]
        elif gps[0] == "SAURON-AB":
            gps = ["50.5580", "-61.0745"]
        elif gps[0] == "EGLAREST-AB":
            gps = ["50.3335", "-60.6990"]
        elif gps[0] == "ANGBAND-AB":
            gps = ["49.0181", "-60.1556"]
        elif gps[0] == "RHUN-AB":
            gps = ["49.0841", "-60.9777"]
        elif gps[0] == "SAURUMAN-AB":
            gps = ["50.0", "-65.7590"]

        alt = self.builder.get_object("ALT").get_text()
        vel = self.builder.get_object("VEL").get_text()
        id = self.builder.get_object("ID").get_text()
        if id == "":
            return
        id = id.zfill(3)
        veh = self.builder.get_object("TYPE").get_active_id() + "-" + id

        if int(alt) >= 0 and int(vel) >= 0:
            subprocess.Popen(
                ["./goto", veh, gps[0], gps[1], alt, vel,
                 self.IP.get_text()])
        return

    def button_gotoland(self, button, data=None):
        gps = self.builder.get_object("GPS").get_text()
        gps = gps.replace(",", " ")
        gps = gps.split()

        if gps[0] == "ELVENKING-AB":
            gps = ["50.6533", "-68.7170"]
        elif gps[0] == "EDORAS-AB":
            gps = ["49.1819", "-68.3639"]
        elif gps[0] == "ETHRING-AB":
            gps = ["49.6165", "-67.7075"]
        elif gps[0] == "EDHELLOND-AB":
            gps = ["50.2223", "-66.2656"]
        elif gps[0] == "SHELOB-AB":
            gps = ["50.5572", "-63.4119"]
        elif gps[0] == "SMAUG-AB":
            gps = ["50.8405", "-62.4290"]
        elif gps[0] == "SCATHA-AB":
            gps = ["50.2306", "-62.1124"]
        elif gps[0] == "SAURON-AB":
            gps = ["50.5580", "-61.0745"]
        elif gps[0] == "EGLAREST-AB":
            gps = ["50.3335", "-60.6990"]
        elif gps[0] == "ANGBAND-AB":
            gps = ["49.0181", "-60.1556"]
        elif gps[0] == "RHUN-AB":
            gps = ["49.0841", "-60.9777"]
        elif gps[0] == "SAURUMAN-AB":
            gps = ["50.0", "-65.7590"]

        alt = self.builder.get_object("ALT").get_text()
        vel = self.builder.get_object("VEL").get_text()
        id = self.builder.get_object("ID").get_text()
        if id == "":
            return
        id = id.zfill(3)
        veh = self.builder.get_object("TYPE").get_active_id() + "-" + id
        if int(alt) >= 0 and int(vel) >= 0:
            subprocess.Popen([
                "./gotoland", veh, gps[0], gps[1], alt, vel,
                self.IP.get_text()
            ])

    def button_launch(self, button, data=None):
        id = self.builder.get_object("ID").get_text()
        if id == "":
            return
        id = id.zfill(3)
        veh = self.builder.get_object("TYPE").get_active_id() + "-" + id
        subprocess.Popen(["./launch", veh, self.IP.get_text()])

    def button_land(self, button, data=None):
        id = self.builder.get_object("ID").get_text()
        if id == "":
            return
        id = id.zfill(3)
        veh = self.builder.get_object("TYPE").get_active_id() + "-" + id
        subprocess.Popen(["./land", veh, self.IP.get_text()])

    def button_status(self, button, data=None):
        id = self.builder.get_object("ID").get_text()
        if id == "":
            return
        id = id.zfill(3)
        veh = self.builder.get_object("TYPE").get_active_id() + "-" + id
        subprocess.Popen(["./status", veh, self.IP.get_text()])

    def button_payload_fire(self, button, data=None):
        action = "1"
        payid = "0"

        if self.builder.get_object("LANCE").get_active():
            payid = "3"
        if self.builder.get_object("CAMERA").get_active():
            payid = "4"
        if self.builder.get_object("PHOSPHEX").get_active():
            payid = "7"
        if self.builder.get_object("PHOSPHEXREM").get_active():
            payid = "8"
        if self.builder.get_object("ANTIMATTER").get_active():
            payid = "10"
        if self.builder.get_object("SEEKER").get_active():
            payid = "15"
        if payid == "0":
            return

        tgt = self.builder.get_object("TGT_ID").get_text()
        if tgt == "":
            tgt = "0"

        id = self.builder.get_object("ID").get_text()
        if id == "":
            return
        id = id.zfill(3)
        veh = self.builder.get_object("TYPE").get_active_id() + "-" + id
        subprocess.Popen(
            ["./executepayload", veh, payid, action, tgt,
             self.IP.get_text()])

    #jettisoning a payload = drone burns less fuel
    def button_payload_jettison(self, button, data=None):
        action = "0"
        payid = "0"

        if self.builder.get_object("LANCE").get_active():
            payid = "3"
        if self.builder.get_object("PHOSPHEX").get_active():
            payid = "7"
        if self.builder.get_object("PHOSPHEXREM").get_active():
            payid = "8"
        if self.builder.get_object("ANTIMATTER").get_active():
            payid = "10"
        if self.builder.get_object("SEEKER").get_active():
            payid = "15"
        if payid == "0":
            return

        tgt = self.builder.get_object("TGT_ID").get_text()
        if tgt == "":
            tgt = "0"

        id = self.builder.get_object("ID").get_text()
        if id == "":
            return
        id = id.zfill(3)
        veh = self.builder.get_object("TYPE").get_active_id() + "-" + id
        subprocess.Popen(
            ["./executepayload", veh, payid, action, tgt,
             self.IP.get_text()])

    def button_loadpayload(self, button, data=None):
        payid = "0"

        if self.builder.get_object("LANCE").get_active():
            payid = "3"
        if self.builder.get_object("FUEL").get_active():
            payid = "5"
        if self.builder.get_object("PHOSPHEX").get_active():
            payid = "7"
        if self.builder.get_object("PHOSPHEXREM").get_active():
            payid = "8"
        if self.builder.get_object("ANTIMATTER").get_active():
            payid = "10"
        if self.builder.get_object("SEEKER").get_active():
            payid = "15"
        if payid == "0":
            return

        id = self.builder.get_object("ID").get_text()
        if id == "":
            return
        id = id.zfill(3)
        veh = self.builder.get_object("TYPE").get_active_id() + "-" + id
        amt = self.builder.get_object("AMT").get_text()
        subprocess.Popen(
            ["./loadpayload", veh, payid, amt,
             self.IP.get_text()])


if __name__ == "__main__":
    main = IcarusUI()
    Gtk.main()
