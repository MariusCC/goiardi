.. _reporting:

Reporting
=========

Goiardi now supports Chef's reporting facilities. Nothing needs to be enabled in goiardi to use this, but changes are required with the client. See http://docs.chef.io/reporting.html for details on how to enable reporting and how to use it.

There is a goiardi extension to reporting: a "status" query parameter may be passed in a GET request that lists reports to limit the reports returned to ones that match the status, so you can read only reports of chef runs that were successful, failed, or started but haven't completed yet. Valid values for the "status" parameter are "started", "success", and "failure".

To use reporting, you'll either need the Chef knife-reporting plugin, or use the knife-goiardi-reporting plugin that supports querying runs by status. It's available on rubygems, or on github at https://github.com/ctdk/knife-goiardi-reporting.
