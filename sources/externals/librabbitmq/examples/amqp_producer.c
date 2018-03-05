/*
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1/GPL 2.0
 *
 * The contents of this file are subject to the Mozilla Public License
 * Version 1.1 (the "License"); you may not use this file except in
 * compliance with the License. You may obtain a copy of the License
 * at http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS"
 * basis, WITHOUT WARRANTY OF ANY KIND, either express or implied. See
 * the License for the specific language governing rights and
 * limitations under the License.
 *
 * The Original Code is librabbitmq.
 *
 * The Initial Developer of the Original Code is VMware, Inc.
 * Portions created by VMware are Copyright (c) 2007-2012 VMware, Inc.
 *
 * Portions created by Tony Garnock-Jones are Copyright (c) 2009-2010
 * VMware, Inc. and Tony Garnock-Jones.
 *
 * All rights reserved.
 *
 * Alternatively, the contents of this file may be used under the terms
 * of the GNU General Public License Version 2 or later (the "GPL"), in
 * which case the provisions of the GPL are applicable instead of those
 * above. If you wish to allow use of your version of this file only
 * under the terms of the GPL, and not to allow others to use your
 * version of this file under the terms of the MPL, indicate your
 * decision by deleting the provisions above and replace them with the
 * notice and other provisions required by the GPL. If you do not
 * delete the provisions above, a recipient may use your version of
 * this file under the terms of any one of the MPL or the GPL.
 *
 * ***** END LICENSE BLOCK *****
 */

#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include <stdint.h>
#include <amqp.h>
#include <amqp_framing.h>

#include "utils.h"

#define SUMMARY_EVERY_US 1000000

static void send_batch(amqp_connection_state_t conn,
		       char const *queue_name,
		       int rate_limit,
		       int message_count)
{
  uint64_t start_time = now_microseconds();
  int i;
  int sent = 0;
  int previous_sent = 0;
  uint64_t previous_report_time = start_time;
  uint64_t next_summary_time = start_time + SUMMARY_EVERY_US;

  char message[256];
  amqp_bytes_t message_bytes;

  for (i = 0; i < sizeof(message); i++) {
    message[i] = i & 0xff;
  }

  message_bytes.len = sizeof(message);
  message_bytes.bytes = message;

  for (i = 0; i < message_count; i++) {
    uint64_t now = now_microseconds();

    die_on_error(amqp_basic_publish(conn,
				    1,
				    amqp_cstring_bytes("amq.direct"),
				    amqp_cstring_bytes(queue_name),
				    0,
				    0,
				    NULL,
				    message_bytes),
		 "Publishing");
    sent++;
    if (now > next_summary_time) {
      int countOverInterval = sent - previous_sent;
      double intervalRate = countOverInterval / ((now - previous_report_time) / 1000000.0);
      printf("%d ms: Sent %d - %d since last report (%d Hz)\n",
	     (int)(now - start_time) / 1000, sent, countOverInterval, (int) intervalRate);

      previous_sent = sent;
      previous_report_time = now;
      next_summary_time += SUMMARY_EVERY_US;
    }

    while (((i * 1000000.0) / (now - start_time)) > rate_limit) {
      microsleep(2000);
      now = now_microseconds();
    }
  }

  {
    uint64_t stop_time = now_microseconds();
    int total_delta = stop_time - start_time;

    printf("PRODUCER - Message count: %d\n", message_count);
    printf("Total time, milliseconds: %d\n", total_delta / 1000);
    printf("Overall messages-per-second: %g\n", (message_count / (total_delta / 1000000.0)));
  }
}

int main(int argc, char const * const *argv) {
  char const *hostname;
  int port;
  int rate_limit;
  int message_count;

  int sockfd;
  amqp_connection_state_t conn;

  if (argc < 5) {
    fprintf(stderr, "Usage: amqp_producer host port rate_limit message_count\n");
    return 1;
  }

  hostname = argv[1];
  port = atoi(argv[2]);
  rate_limit = atoi(argv[3]);
  message_count = atoi(argv[4]);

  conn = amqp_new_connection();

  die_on_error(sockfd = amqp_open_socket(hostname, port), "Opening socket");
  amqp_set_sockfd(conn, sockfd);
  die_on_amqp_error(amqp_login(conn, "/", 0, 131072, 0, AMQP_SASL_METHOD_PLAIN, "guest", "guest"),
		    "Logging in");
  amqp_channel_open(conn, 1);
  die_on_amqp_error(amqp_get_rpc_reply(conn), "Opening channel");

  send_batch(conn, "test queue", rate_limit, message_count);

  die_on_amqp_error(amqp_channel_close(conn, 1, AMQP_REPLY_SUCCESS), "Closing channel");
  die_on_amqp_error(amqp_connection_close(conn, AMQP_REPLY_SUCCESS), "Closing connection");
  die_on_error(amqp_destroy_connection(conn), "Ending connection");
  return 0;
}
