'use strict';
var dataProvider = require('../data/skills.js');
/**
 * Operations on /skills
 */
module.exports = {
    /**
     * summary: Retrieves skills.
     * description: 
     * parameters: 
     * produces: application/json, text/json
     * responses: 200
     */
    get: function skills_get(req, res, next) {
        /**
         * Get the data for response 200
         * For response `default` status 200 is used.
         */
        var status = 200;
        var provider = dataProvider['get']['200'];
        provider(req, res, function (err, data) {
            if (err) {
                next(err);
                return;
            }
            res.status(status).send(data && data.responses);
        });
    }
};
