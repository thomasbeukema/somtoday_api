package com.beech32.myschool;

import android.os.AsyncTask;

import org.apache.http.HttpEntity;
import org.apache.http.HttpResponse;
import org.apache.http.client.HttpClient;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.DefaultHttpClient;
import org.apache.http.params.BasicHttpParams;
import org.apache.http.params.DefaultedHttpParams;
import org.apache.http.params.HttpConnectionParams;
import org.apache.http.params.HttpParams;

import java.io.BufferedReader;
import java.io.InputStreamReader;

public class CallAPI extends AsyncTask<String, String, String> {

    @Override
    protected String doInBackground(String... params) {

        String urlString = params[0].replace("%0A", "");

        String resultToDisplay = "";

        BufferedReader in = null;
        StringBuilder sb = null;

        try {
            HttpClient httpclient = new DefaultHttpClient();
            HttpGet get = new HttpGet(urlString);
            HttpParams httpParams = new BasicHttpParams();
            HttpConnectionParams.setConnectionTimeout(httpParams, 3000);
            HttpConnectionParams.setSoTimeout(httpParams, 5000);
            get.setParams(httpParams);
            HttpResponse response = httpclient.execute(get);

            HttpEntity entity = response.getEntity();

            in = new BufferedReader(new InputStreamReader(entity.getContent(), "UTF-8"));
            sb = new StringBuilder();

            String inputStr;

            while((inputStr = in.readLine()) != null) {
                sb.append(inputStr);
            }



        } catch (Exception e) {
            System.out.println(e);

            return "OFFLINE";
        }

        resultToDisplay = sb.toString();

        return resultToDisplay;

    }

    protected void onPostExecute(String result) {

    }

}
