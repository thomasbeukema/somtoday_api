import android.util.Base64;

import org.json.JSONArray;
import org.json.JSONObject;

import java.io.UnsupportedEncodingException;
import java.net.URLEncoder;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

import java.util.ArrayList;

public class Somtoday {

    private String username;
    private String password;
    private String brin;
    private String school;
    private String fullName;
    private String pupilID;
    private boolean rightCredentials;
    private String baseURL;

    public Somtoday() {
    }

    public void setCredentials(String userName, String pass, String BRIN, String schoolname) {
        username = userName;
        password = pass;
        brin = BRIN;
        school = schoolname;
    }

    public String getItem(String item) {
        switch (item) {
            case "username":
                return username;
            case "password":
                return password;
            case "brin":
                return brin;
            case "school":
                return school;
            case "fullName":
                return fullName;
            case "pupilID":
                return pupilID;
            default:
                return "NAM";
        }
    }

    public boolean login(ArrayList<String> credentials) {
        try {
            String url = "http://somtoday.nl/" + URLEncoder.encode(credentials.get(3), "ASCII") + "/services/mobile/v10/Login/CheckMultiLoginB64/" + credentials.get(0) + "/" + URLEncoder.encode(credentials.get(1), "ASCII") + "/" + URLEncoder.encode(credentials.get(2), "ASCII");
            String asyncResult = new CallAPI().execute(url, "login").get();

            JSONObject jsonObject = new JSONObject(asyncResult);
            if (jsonObject.getString("error").equals("SUCCESS")) {

                rightCredentials = true;

                JSONArray tempJsonObj = jsonObject.getJSONArray("leerlingen");
                fullName = tempJsonObj.getJSONObject(0).getString("fullName");
                pupilID = tempJsonObj.getJSONObject(0).getString("leerlingId");

                baseURL = "https://somtoday.nl/" + school + "/services/mobile/v10/";

                return true;

            } else {
                return false;
            }

        } catch (Exception e) {
            return false;
        }
    }

    public JSONObject getRooster(int daysAhead, ArrayList<String> credentials) {
        long date = (System.currentTimeMillis() + (daysAhead * 86400000));

        try {
            String url = baseURL + "Agenda/GetMultiStudentAgendaB64/" + URLEncoder.encode(credentials.get(0), "ASCII") + "/" + URLEncoder.encode(credentials.get(1), "ASCII") + "/" + URLEncoder.encode(credentials.get(2), "ASCII") + "/" + String.valueOf(date) + "/" + URLEncoder.encode(pupilID, "ASCII");
            String asyncResult = new CallAPI().execute(url).get();

            return new JSONObject(asyncResult);
        } catch (Exception e) {

        }

        return new JSONObject();
    }

    public JSONObject getHuiswerk(int daysAhead, ArrayList<String> credentials) {
        try {
            String url = baseURL + "Agenda/GetMultiStudentAgendaHuiswerkMetMaxB64/" + URLEncoder.encode(credentials.get(0), "ASCII") + "/" + URLEncoder.encode(credentials.get(1), "ASCII") + "/" + URLEncoder.encode(credentials.get(2), "ASCII") + "/" + String.valueOf(daysAhead) + "/" + URLEncoder.encode(pupilID, "ASCII");
            String asyncResult = new CallAPI().execute(url).get();
            return new JSONObject(asyncResult);
        } catch (Exception e) {

        }

        return new JSONObject();
    }

    public JSONObject getCijfers(int state, ArrayList<String> credentials) {
        if (state == 0) {
            try {
                String url = baseURL + "Cijfers/GetMultiCijfersRecentB64/" + URLEncoder.encode(credentials.get(0), "ASCII") + "/" + URLEncoder.encode(credentials.get(1), "ASCII") + "/" + URLEncoder.encode(credentials.get(2), "ASCII") + "/" + URLEncoder.encode(pupilID, "ASCII");
                String asyncResult = new CallAPI().execute(url).get();
                return new JSONObject(asyncResult);
            } catch (Exception e) {

            }
        } else {
            try {
                String url = baseURL + "Cijfers/GetCijfersRecentMetMaxB64/" + URLEncoder.encode(credentials.get(0), "ASCII") + "/" + URLEncoder.encode(credentials.get(1), "ASCII") + "/" + URLEncoder.encode(credentials.get(2), "ASCII") + "/1000";
                String asyncResult = new CallAPI().execute(url).get();
                return new JSONObject(asyncResult);
            } catch (Exception e) {

            }
        }

        return new JSONObject();
    }

    public boolean changeHwStatus(String afspraakID, String hwID, boolean done, ArrayList<String> credentials) {
        try {
            String status = "0";

            if (done) {
                status = "1";
            }

            String url = baseURL + "Agenda/HuiswerkDoneB64/" + URLEncoder.encode(credentials.get(0), "ASCII") + "/" + URLEncoder.encode(credentials.get(1), "ASCII") + "/" + URLEncoder.encode(credentials.get(2), "ASCII") + "/" + afspraakID + "/" + hwID + "/" + status;
            String asyncResult = new CallAPI().execute(url).get();
        } catch (Exception e) {

    }

        return true;
    }

    private byte[] sha1(String input) throws NoSuchAlgorithmException, UnsupportedEncodingException {
        MessageDigest sha1 = MessageDigest.getInstance("SHA1");
        byte[] data = input.getBytes("ASCII");
        byte[] hash = sha1.digest(data);

        return hash;
    }

    public static String bytesToHex(byte[] bytes) {
        final char[] hexArray = "0123456789ABCDEF".toCharArray();
        char[] hexChars = new char[bytes.length * 2];
        for ( int j = 0; j < bytes.length; j++ ) {
            int v = bytes[j] & 0xFF;
            hexChars[j * 2] = hexArray[v >>> 4];
            hexChars[j * 2 + 1] = hexArray[v & 0x0F];
        }
        return new String(hexChars);
    }

    public ArrayList<String> getUrlParams() {

        String userName = null;
        byte[] passwordEncoded = null;
        String passwordEncodedFinal = null;

        try {
            userName = Base64.encodeToString(username.getBytes(), Base64.NO_WRAP);
            passwordEncoded = sha1(password);
            passwordEncoded = Base64.encode(passwordEncoded, Base64.NO_WRAP);
            passwordEncodedFinal = bytesToHex(passwordEncoded);
        } catch (Exception e) {

        }

        ArrayList<String> finalList = new ArrayList<>();
        finalList.add(0, userName);
        finalList.add(1, passwordEncodedFinal);
        finalList.add(2, brin);
        finalList.add(3, school);

        return finalList;
    }

}
