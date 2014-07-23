Send Event With Perl
====================

Requirements
------------

Install Net::RabbitFoot
~~~~~~~~~~~~~~~~~~~~~~~~

``cpan -i Net::RabbitFoot``

Install JSON::XS
~~~~~~~~~~~~~~~~~

``cpan -i JSON::XS``

Example
-------

.. code-block:: perl

	

    #!/usr/bin/perl -w

    use strict;
    use warnings;
    use Net::RabbitFoot;
    use JSON::XS;

    my $connector_name = "rrd";
    my $component = "MACHINE_TEST35";
    my $resource = "RESOURCE_TEST";
    my $output = "Refreshing Resource Perf Data";
    my $timestamp = time;
    my $perf_data_metric = "used";
    my $perf_data_unit = "Hz";
    my $perf_data_value = 10.8;
    my $perf_data_type = "GAUGE";

    my $hash = {
        "connector_name" => $connector_name,
        "event_type" => "log",
        "source_type" => $resource,
        "component" => $component,
        "resource" => $resource,
        "state" => 0,
        "state_type" => 1,
        "output" => $output,
        "timestamp" => $timestamp,
        "perf_data_array" => [{
            "metric" => $perf_data_metric, 
            "unit" => $perf_data_unit, 
            "value" => $perf_data_value,
            "type" => $perf_data_type
            }]
        };

    my $json = JSON::XS->new->utf8->space_after->encode ($hash);

    my $connect = Net::RabbitFoot->new()->load_xml_spec()->connect(
        host => '127.0.0.1',
        port => 5672,
        user => 'guest',
        pass => 'guest',
        vhost => 'canopsis',
    );

    my $channel = $connect->open_channel();

    ###To send many events, put your loops here ...
    $channel->publish(
        exchange => 'canopsis.events',
        routing_key => "cli.".$connector_name.".log.".$resource.".".$component.".".$resource,
        body => $json,
    );
    ### End of loop ! (rebuild your json in this loop if you want to use it)
    $connect->close();
  
