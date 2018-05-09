package main

import (
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
)

type stat struct {
	name   string
	stat   string
	regexp *regexp.Regexp
	value  prometheus.ValueType
	help   string
	desc   *prometheus.Desc
}

var opensipsStats = map[string][]stat{
	"core": {
		{
			name:  "received_requests_total",
			stat:  "rcv_requests",
			value: prometheus.CounterValue,
			help:  "The total number of received requests by OpenSIPS",
		},
		{
			name:  "received_replies_total",
			stat:  "rcv_replies",
			value: prometheus.CounterValue,
			help:  "The total number of received replies by OpenSIPS",
		},
		{
			name:  "forwarded_requests_total",
			stat:  "fwd_requests",
			value: prometheus.CounterValue,
			help:  "Total number of stateless forwarded requests by OpenSIPS",
		},
		{
			name:  "forwarded_replies_total",
			stat:  "fwd_replies",
			value: prometheus.CounterValue,
			help:  "Total number of stateless forwarded replies by OpenSIPS",
		},
		{
			name:  "dropped_requests_total",
			stat:  "drop_requests",
			value: prometheus.CounterValue,
			help:  "Total number of requests dropped even before entering the script routing logic",
		},
		{
			name:  "dropped_replies_total",
			stat:  "drop_replies",
			value: prometheus.CounterValue,
			help:  "Total number of replies dropped even before entering the script routing logic",
		},
		{
			name:  "error_requests_total",
			stat:  "err_requests",
			value: prometheus.CounterValue,
			help:  "Total number of bogus or invalid requests",
		},
		{
			name:  "error_replies_total",
			stat:  "err_replies",
			value: prometheus.CounterValue,
			help:  "Total number of bogus or invalid replies",
		},
		{
			name:  "bad_uris_received_total",
			stat:  "bad_URIs_rcvd",
			value: prometheus.CounterValue,
			help:  "Total number of URIs that OpenSIPS failed to parse",
		},
		{
			name:  "unsupported_methods_total",
			stat:  "unsupported_methods",
			value: prometheus.CounterValue,
			help:  "Total number of non-standard methods encountered by OpenSIPS while parsing SIP methods",
		},
		{
			name:  "bad_message_headers_total",
			stat:  "bad_msg_hdr",
			value: prometheus.CounterValue,
			help:  "Total number of SIP headers that OpenSIPS failed to parse",
		},
		{
			name:  "uptime_seconds_total",
			stat:  "timestamp",
			value: prometheus.CounterValue,
			help:  "The number of seconds elapsed from OpenSIPS starting",
		},
	},
	"dialog": {
		{
			name:  "active_dialogs",
			stat:  "active_dialogs",
			value: prometheus.GaugeValue,
			help:  "Number of active dialogs",
		},
		{
			name:  "early_dialogs",
			stat:  "early_dialogs",
			value: prometheus.GaugeValue,
			help:  "Number of early dialogs",
		},
		{
			name:  "processed_dialogs_total",
			stat:  "processed_dialogs",
			value: prometheus.CounterValue,
			help:  "Total number of processed dialogs",
		},
		{
			name:  "expired_dialogs_total",
			stat:  "expired_dialogs",
			value: prometheus.CounterValue,
			help:  "Total number of expired dialogs",
		},
		{
			name:  "failed_dialogs_total",
			stat:  "failed_dialogs",
			value: prometheus.CounterValue,
			help:  "Total number of failed dialogs",
		},
		{
			name:   "replication_messages_sent_total",
			regexp: regexp.MustCompile(`^(?P<operation>.+)_sent$`),
			value:  prometheus.CounterValue,
			help:   "Total number of replication messages sent",
		},
		{
			name:   "replication_messages_received_total",
			regexp: regexp.MustCompile(`^(?P<operation>.+)_recv$`),
			value:  prometheus.CounterValue,
			help:   "Total number of replication messages received",
		},
	},
	"load": {
		{
			name:  "load",
			stat:  "load",
			value: prometheus.GaugeValue,
			help:  "The real time load of core OpenSIPS processes",
		},
		{
			name:  "load_all",
			stat:  "load-all",
			value: prometheus.GaugeValue,
			help:  "The real time load of all OpenSIPS processes",
		},
		{
			name:   "process_load",
			regexp: regexp.MustCompile(`^load-proc-(?P<id>\d+)$`),
			value:  prometheus.GaugeValue,
			help:   "he real time load of the OpenSIPS process #id",
		},
	},
	"msilo": {
		{
			name:  "stored_messages_total",
			stat:  "stored_messages",
			value: prometheus.CounterValue,
			help:  "Total number of stored messages",
		},
		{
			name:  "dumped_messages_total",
			stat:  "dumped_messages",
			value: prometheus.CounterValue,
			help:  "Total number of dumped messages",
		},
		{
			name:  "failed_messages_total",
			stat:  "failed_messages",
			value: prometheus.CounterValue,
			help:  "Total number of failed messages",
		},
		{
			name:  "dumped_reminders_total",
			stat:  "dumped_reminders",
			value: prometheus.CounterValue,
			help:  "Total number of dumped reminders",
		},
		{
			name:  "failed_reminders_total",
			stat:  "failed_reminders",
			value: prometheus.CounterValue,
			help:  "Total number of failed reminders",
		},
	},
	"nat_traversal": {
		{
			name:  "keepalive_endpoints",
			stat:  "keepalive_endpoints",
			value: prometheus.GaugeValue,
			help:  "Current number of keepalive endpoints",
		},
		{
			name:  "registered_endpoints",
			stat:  "registered_endpoints",
			value: prometheus.GaugeValue,
			help:  "Current number of registered endpoints",
		},
		{
			name:  "subscribed_endpoints",
			stat:  "subscribed_endpoints",
			value: prometheus.GaugeValue,
			help:  "Current number of subscribed endpoints",
		},
		{
			name:  "dialog_endpoints",
			stat:  "dialog_endpoints",
			value: prometheus.GaugeValue,
			help:  "Current number of dialog endpoints",
		},
	},
	"net": {
		{
			name:   "waiting_bytes",
			regexp: regexp.MustCompile(`^waiting_(?P<transport>.+)$`),
			value:  prometheus.GaugeValue,
			help:   "The number of bytes waiting to be consumed on interfaces that OpenSIPS is listening on",
		},
	},
	"pkmem": {
		{
			name:   "total_size_bytes",
			regexp: regexp.MustCompile(`^(?P<id>\d+)-total_size$`),
			value:  prometheus.GaugeValue,
			help:   "The total size of private memory available to OpenSIPS process #id",
		},
		{
			name:   "used_size_bytes",
			regexp: regexp.MustCompile(`^(?P<id>\d+)-used_size$`),
			value:  prometheus.GaugeValue,
			help:   "The total size of private memory used by OpenSIPS process #id",
		},
		{
			name:   "real_used_size_bytes",
			regexp: regexp.MustCompile(`^(?P<id>\d+)-real_used_size$`),
			value:  prometheus.GaugeValue,
			help:   "The total size of private memory (including overhead) used by OpenSIPS process #id",
		},
		{
			name:   "max_used_size_bytes",
			regexp: regexp.MustCompile(`^(?P<id>\d+)-max_used_size$`),
			value:  prometheus.GaugeValue,
			help:   "The maximum amount of private memory ever used by OpenSIPS process #id",
		},
		{
			name:   "free_size_bytes",
			regexp: regexp.MustCompile(`^(?P<id>\d+)-free_size$`),
			value:  prometheus.GaugeValue,
			help:   "The free private memory available for OpenSIPS process #id",
		},
		{
			name:   "fragments",
			regexp: regexp.MustCompile(`^(?P<id>\d+)-fragments$`),
			value:  prometheus.GaugeValue,
			help:   "The number of fragments in the private memory for OpenSIPS process #",
		},
	},
	"registrar": {
		{
			name:  "max_expires",
			stat:  "max_expires",
			value: prometheus.GaugeValue,
			help:  "The value of the max_expires module parameter",
		},
		{
			name:  "max_contacts",
			stat:  "max_contacts",
			value: prometheus.GaugeValue,
			help:  "The value of the max_contacts module parameter",
		},
		{
			name:  "default_expires",
			stat:  "default_expire",
			value: prometheus.GaugeValue,
			help:  "The value of the default_expires module parameter",
		},
		{
			name:  "accepted_registrations_total",
			stat:  "accepted_registrations",
			value: prometheus.CounterValue,
			help:  "Total number of accepted registrations",
		},
		{
			name:  "rejected_registrations_total",
			stat:  "rejected_registrations",
			value: prometheus.CounterValue,
			help:  "Total number of rejected registrations",
		},
	},
	"shmem": {
		{
			name:  "total_size_bytes",
			stat:  "total_size",
			value: prometheus.GaugeValue,
			help:  "The total size of shared memory available to OpenSIPS processes",
		},
		{
			name:  "used_size_bytes",
			stat:  "used_size",
			value: prometheus.GaugeValue,
			help:  "The total size of shared memory used by OpenSIPS processes",
		},
		{
			name:  "real_used_size_bytes",
			stat:  "real_used_size",
			value: prometheus.GaugeValue,
			help:  "The total size of shared memory used (including overhead) by OpenSIPS processes",
		},
		{
			name:  "max_used_size_bytes",
			stat:  "max_used_size",
			value: prometheus.GaugeValue,
			help:  "The maximum amount of shared memory used by OpenSIPS processes",
		},
		{
			name:  "free_size_bytes",
			stat:  "free_size",
			value: prometheus.GaugeValue,
			help:  "The amount of free shared memory available to OpenSIPS processes",
		},
		{
			name:  "fragments",
			stat:  "fragments",
			value: prometheus.GaugeValue,
			help:  "The number of fragments in the shared memory used by OpenSIPS processes",
		},
	},
	"sipcapture": {
		{
			name:  "captured_requests_total",
			stat:  "captured_requests",
			value: prometheus.CounterValue,
			help:  "Total number of SIP requests captured",
		},
		{
			name:  "captured_replies_total",
			stat:  "captured_replies",
			value: prometheus.CounterValue,
			help:  "Total number of SIP replies captured",
		},
	},
	"siptrace": {
		{
			name:  "traced_requests_total",
			stat:  "traced_requests",
			value: prometheus.CounterValue,
			help:  "Total number of traced requests",
		},
		{
			name:  "traced_replies_total",
			stat:  "traced_replies",
			value: prometheus.CounterValue,
			help:  "Total number of traced replies",
		},
	},
	"sl": {
		{
			name:   "sent_replies",
			regexp: regexp.MustCompile(`^(?P<code>[2-6]xx)_replies$`),
			value:  prometheus.CounterValue,
			help:   "Total number of sent replies by status code",
		},
		{
			name:  "sent_replies_total",
			stat:  "sent_replies",
			value: prometheus.CounterValue,
			help:  "Total number of sent replies",
		},
		{
			name:  "sent_error_replies_total",
			stat:  "sent_err_replies",
			value: prometheus.CounterValue,
			help:  "Total number of sent error replies",
		},
		{
			name:  "received_acks_total",
			stat:  "received_ACKs",
			value: prometheus.CounterValue,
			help:  "Total number of ACK replies received",
		},
	},
	"sst": {
		{
			name:  "expired_sst_total",
			stat:  "expired_sst",
			value: prometheus.CounterValue,
			help:  "Total number of expired SST sessions",
		},
	},
	"tm": {
		{
			name:  "received_replies_total",
			stat:  "received_replies",
			value: prometheus.CounterValue,
			help:  "Total number of replies received",
		},
		{
			name:  "relayed_replies_total",
			stat:  "relayed_replies",
			value: prometheus.CounterValue,
			help:  "Total number of replies relayed",
		},
		{
			name:  "local_replies_total",
			stat:  "local_replies",
			value: prometheus.CounterValue,
			help:  "Total number of local replies sent",
		},
		{
			name:  "uas_transactions_total",
			stat:  "UAS_transactions",
			value: prometheus.CounterValue,
			help:  "Total number of UAS transactions",
		},
		{
			name:  "uac_transactions_total",
			stat:  "UAC_transactions",
			value: prometheus.CounterValue,
			help:  "Total number of UAC transactions",
		},
		{
			name:   "transactions_total",
			regexp: regexp.MustCompile(`^(?P<code>[2-6]xx)_transactions$`),
			value:  prometheus.CounterValue,
			help:   "Total number of transactions by status code",
		},
		{
			name:  "inuse_transactions",
			stat:  "inuse_transactions",
			value: prometheus.GaugeValue,
			help:  "Number of transactions currently in-use",
		},
	},
	"uri": {
		{
			name:  "positive_checks_total",
			stat:  "positive_checks",
			value: prometheus.CounterValue,
			help:  "Total number of positive URI checks",
		},
		{
			name:  "negative_checks_total",
			stat:  "negative_checks",
			value: prometheus.CounterValue,
			help:  "Total number of negative URI checks",
		},
	},
	"usrloc": {
		{
			name:  "registered_users",
			stat:  "registered_users",
			value: prometheus.GaugeValue,
			help:  "Current number of registered users",
		},
	},
}

func init() {
	for subsys, stats := range opensipsStats {
		for i, stat := range stats {
			if stat.desc == nil {
				labels := []string{}
				if stat.regexp != nil {
					labels = stat.regexp.SubexpNames()[1:]
				}

				stats[i].desc = prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsys, stat.name),
					stat.help,
					labels,
					nil,
				)
			}
		}
	}
}
