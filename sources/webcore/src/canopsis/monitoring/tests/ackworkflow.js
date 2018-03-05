// module('Ack workflow', {
//     setup: function () {
//         App.reset();
//     },
//     teardown: function() {
//         $.mockjax.clear();
//     }
// });

// test('Test ack workflow', function() {

//     expect(4);

//     visit('/userview/view.event');

//     andThen(function() {
//         var listWidgetsCount = find('.widget.list').length;
//         equal(listWidgetsCount, 1, 'One list widget found');
//     });

//     andThen(function() {
//         var json = {
//             total:0,
//             data:[{
//                 "status":1,
//                 "crecord_type":"event",
//                 "event_type":"check",
//                 "timestamp":1438698464,
//                 "component":"A",
//                 "source_type":"resource",
//                 "id":"Engine.engine.check.resource.A.B",
//                 "resource":"B",
//                 "event_id":"Engine.engine.check.resource.A.B",
//                 "connector":"Engine",
//                 "state":2,
//                 "connector_name":"engine",
//                 "output":"",
//                 "_id":"Engine.engine.check.resource.A.B",
//                 "rk":"Engine.engine.check.resource.A.B"
//             }],
//             success:true
//         };

//         stubEndpointForHttpRequest('/rest/events', json);

//         //refresh data
//         click('.canopsis-toolbar .glyphicon-refresh');

//     });

//     andThen(function() {
//         click('.listline button');
//     });

//     andThen(function() {
//         fillIn('.modal-content .ember-text-field', 'ticketNumber');
//         fillIn('.modal-content .ember-text-area', 'reason');
//     });

//     andThen(function() {
//         stubEndpointForHttpRequest(
//             '/event',
//             {success: true},
//             function (settings) {
//                 sendEvent = JSON.parse(settings.data.event)[0];
//                 equal(sendEvent.component, 'A', 'Expect the event component is equal to "A"');
//                 equal(sendEvent.output, 'reason', 'Expect the event ouptut is equal to "reason"');
//                 equal(sendEvent.event_type, 'ack', 'Expect the event type is equal to "ack"');
//             }
//         );

//         click('.modal-footer .btn-success');
//     });

// });

// test('Test cancel workflow', function() {


//     expect(3);

//     visit('/userview/view.event');

//     andThen(function() {
//         var json = {
//             total:0,
//             data:[{
//                 "status":1,
//                 "crecord_type":"event",
//                 "event_type":"check",
//                 "timestamp":1438698464,
//                 "component":"A",
//                 "source_type":"resource",
//                 "id":"Engine.engine.check.resource.A.B",
//                 "resource":"B",
//                 "event_id":"Engine.engine.check.resource.A.B",
//                 "connector":"Engine",
//                 "state":2,
//                 "connector_name":"engine",
//                 "output":"",
//                 "last_state_change" : 1437984466,
//                 "ack": {
//                     "author": "root",
//                     "comment": "output",
//                     "isAck": true,
//                     "rk": "Engine.engine.check.resource.A.B",
//                     "timestamp": 1437984466
//                 },
//                 "_id":"Engine.engine.check.resource.A.B",
//                 "rk":"Engine.engine.check.resource.A.B"
//             }],
//             success:true
//         };

//         stubEndpointForHttpRequest('/rest/events', json);

//         //refresh data
//         click('.canopsis-toolbar .glyphicon-refresh');

//     });

//     andThen(function() {
//         click('.listline .glyphicon-ban-circle');
//     });

//     andThen(function() {
//         fillIn(".modal-content [name='output']", 'cancel ack');
//     });

//     andThen(function() {
//         stubEndpointForHttpRequest(
//             '/event',
//             {success: true},
//             function (settings) {
//                 sendEvent = JSON.parse(settings.data.event)[0];
//                 equal(sendEvent.component, 'A', 'Expect the event component is equal to "A"');
//                 equal(sendEvent.output, 'cancel ack', 'Expect the event ouptut is equal to "reason"');
//                 equal(sendEvent.event_type, 'ackremove', 'Expect the event type is equal to "ack"');
//             }
//         );

//         click('.modal-footer .btn-primary');
//     });
// });
