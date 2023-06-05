BEGIN;

DROP TABLE IF EXISTS cps_event;

DROP TABLE IF EXISTS fifo_queue;
DROP TABLE IF EXISTS fifo_event;

DROP TABLE IF EXISTS che_event;
DROP TABLE IF EXISTS che_infos;

DROP TABLE IF EXISTS axe_event;
DROP TABLE IF EXISTS axe_periodical;

DROP TABLE IF EXISTS pbehavior_periodical;

DROP TABLE IF EXISTS correlation_event;

DROP TABLE IF EXISTS service_event;

DROP TABLE IF EXISTS dynamic_infos_event;

DROP TABLE IF EXISTS action_event;

DROP TABLE IF EXISTS api_requests;

COMMIT;
