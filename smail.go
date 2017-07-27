/*
Copyright 2017 Mario Kleinsasser and Bernhard Rausch

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"io/ioutil"
	"net"
	"os"

	log "github.com/Sirupsen/logrus"
	gomail "gopkg.in/gomail.v2"
)

// data type for email
type maildata struct {
	From       string
	Subject    string
	To         []string
	Attachment string
	Mx         string
}

// multivalue to address flag
type toarrayFlags []string

func (i *toarrayFlags) String() string {
	return "my string representation"
}

func (i *toarrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func parseFlags() (maildata, error) {

	var toFlags toarrayFlags
	var md maildata

	log.Debug("Parsing command line flags.")
	subject := flag.String("s", "", "Subject text")
	flag.Var(&toFlags, "t", "To address(es)")
	from := flag.String("f", "", "From address")
	attachment := flag.String("a", "", "Attachment [optional]")
	mx := flag.String("m", "", "MX DNS record")
	debug := flag.Bool("d", false, "Debug [default=false]")

	flag.Parse()

	if *debug == true {
		log.SetLevel(log.DebugLevel)
	}

	if *subject == "" {
		return md, errors.New("No subject given!")
	}

	if toFlags == nil {
		return md, errors.New("No to given!")
	}

	if *from == "" {
		return md, errors.New("No subject given!")
	}

	if *mx == "" {
		return md, errors.New("No MX DNS record given")
	}

	md.Subject = *subject
	md.To = toFlags
	md.From = *from
	md.Attachment = *attachment
	md.Mx = *mx

	return md, nil

}

func sendmail(md maildata) error {
	log.Info("Start mailing sequence.")

	var attachmenttext string

	if md.Attachment != "" {
		dat, err := ioutil.ReadFile(md.Attachment)
		if err != nil {
			return err
		}
		attachmenttext = string(dat[:])
		log.Debug(attachmenttext)
	}

	// resolv the mx record
	resolved, err := resolveMx(md)
	if err != nil {
		return nil
	}

	for _, e := range resolved {
		for _, t := range md.To {
			err := send(md.From, t, md.Subject, &attachmenttext, e.Host)
			if err != nil {
				log.Warnf("Cant send email via: %s", e.Host)
				continue
			}
		}
		log.Info("Stopping mailing sequence.")
		return nil
	}

	return errors.New("Cannot send mail via any smtpserver")
}

func send(from string, to string, subject string, body *string, smtpserver string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", *body)

	d := &gomail.Dialer{Host: smtpserver, Port: 25}
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Panic(err)
	}

	log.Infof("Mail send to: %s", to)

	return nil
}

func resolveMx(md maildata) ([]*net.MX, error) {

	resolved, err := net.LookupMX(md.Mx)
	if err != nil {
		log.Info(err)
		return resolved, err
	}

	return resolved, nil
}

func main() {

	// configure logrus logger
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	customFormatter.ForceColors = true
	log.SetFormatter(customFormatter)
	log.SetOutput(os.Stdout)

	// parse the command line flags
	md, err := parseFlags()
	if err != nil {
		log.Panic(err)
	}

	log.Debugf("Given email parameters: %s", md)

	err = sendmail(md)
	if err != nil {
		log.Info(err)
	}

}
