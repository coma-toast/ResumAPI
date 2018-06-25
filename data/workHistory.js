'use strict';
var Mockgen = require('./mockgen.js');
/**
 * Operations on /workHistory
 */
module.exports = {
    /**
     * summary: Retrieves work history.
     * description: 
     * parameters: 
     * produces: application/json, text/json
     * responses: 200
     * operationId: workHistory_get
     */
    get: {
        200: function (req, res, callback) {
            /**
             * Using mock data generator module.
             * Replace this by actual data for the api.
             */
            Mockgen().responses({
                path: '/workHistory',
                operation: 'get',
                response: '200'
            }, callback);
        }
    }
};
