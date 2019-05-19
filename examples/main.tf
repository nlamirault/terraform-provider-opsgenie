# Copyright (C) 2018-2019 Nicolas Lamirault <nicolas.lamirault@gmail.com>

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

terraform {
  required_version = ">= 0.11.0"
}

provider "opsgenie" {
  api_key = "${var.api_key}"
}

resource "opsgenie_user" "first" {
  username  = "john@doe.com"
  full_name = "John Doe"
  role      = "admin"
  locale    = "fr_FR"
  timezone  = "Europe/Paris"
}

resource "opsgenie_user" "second" {
  username  = "jane@doe.com"
  full_name = "Jane Doe"
  role      = "user"
  locale    = "fr_FR"
  timezone  = "Europe/Paris"
}

resource "opsgenie_team" "team_test" {
  name        = "Test"
  description = "This team deals with all the things"

  member {
    username = "${opsgenie_user.first.username}"
    role     = "admin"
  }

  member {
    username = "${opsgenie_user.second.username}"
    role     = "user"
  }
}

resource "opsgenie_contact" "first_contact_email" {
  username = "${opsgenie_user.first.username}"
  to       = "john.doe@doe.com"
  method   = "email"
}

resource "opsgenie_contact" "first_contact_sms" {
  username = "${opsgenie_user.first.username}"
  to       = "33-600000000"
  method   = "sms"
}

resource "opsgenie_schedule" "schedule_doe_ops" {
  name        = "Doe team schedule"
  description = "A schedule for the Doe team"
  owner       = "${opsgenie_team.team_test.name}"
  timezone    = "Europe/Paris"

  rotation = {
    name       = "First rotation"
    start_date = "2019-05-20T08:00:00Z"
    end_date   = "2019-05-24T19:00:00Z"
    type       = "daily"

    participant {
      type     = "user"
      username = "${opsgenie_user.first.username}"
    }
  }

  rotation = {
    name       = "Second rotation"
    start_date = "2019-05-27T08:00:00Z"
    end_date   = "2019-05-31T19:00:00Z"
    type       = "daily"

    participant {
      type     = "user"
      username = "${opsgenie_user.second.username}"
    }
  }
}
