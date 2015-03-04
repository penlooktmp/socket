/**
 * Penlook Project
 *
 * Copyright (c) 2015 Penlook Development Team
 *
 * --------------------------------------------------------------------
 *
 * This program is free software: you can redistribute it and/or
 * modify it under the terms of the GNU Affero General Public License
 * as published by the Free Software Foundation, either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public
 * License along with this program.
 * If not, see <http://www.gnu.org/licenses/>.
 *
 * --------------------------------------------------------------------
 *
 * Authors:
 *     Loi Nguyen       <loint@penlook.com>
 *     Nam Vo           <namvh@penlook.com>
 */

part of socket;

class Option {

    /**
     * HTTP Method
     *
     * @var string
     */
    String method;
    String get Method => method;
    set Method(String method_) => method = method_;

    /**
     * URL
     *
     * @var string
     */
    String url;
    String get Url => url;
    set Url(String url_) => url = url_;

    /**
     * Data package
     *
     * @var string
     */
    String data;
    String get Data => data;
    set Data(String data_) => data = data_;

    /**
     * Request timeout
     *
     * @var int
     */
    int timeout;
    int get Timeout => timeout;
    set Timeout(int timeout_) => timeout = timeout_;

    /**
     * Asynchronous request
     *
     * @var bool
     */
    bool async;
    bool get Async => async;
    set Async(bool async_) => async = async_;

    /**
     * Option constructor
     *
     * @param string method default "GET"
     * @param string url default ""
     * @param string data default "{}"
     * @param int timeout default 60
     * @param bool async default true
     */
    Option({
        String method: "GET",
        String url : "",
        String data : "{}",
        int timeout : 60,
        bool async : true
    }) {
        this.method = method;
        this.url = url;
        this.data = data;
        this.timeout = timeout;
        this.async = async;
    }
}