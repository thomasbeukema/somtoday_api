import hashlib, base64, requests, time, codecs

class FEATURE_NOT_ACTIVATED(Exception):
    pass

class FAILED_AUTHENTICATION(Exception):
    pass

class FAILED_OTHER_TYPE(Exception):
    pass

class SOMTODAY_ERROR_UNKNOWN(Exception):
    pass

class Somtoday:
    leerlingID = None
    fullName = ""

    def __init__(self, credentials):

        """
            Create Somtoday object
            
            Parameters:
            
                credentials:
                    tuple:
                        0 => username
                        1 => password
                        2 => school
                        3 => brin
             
        """
        
        username = bytes(credentials[0], "UTF-8") # Extract username from tuple
        username = base64.b64encode(username) # Encode username
        username = str(username.decode())

        password = hashlib.sha1(credentials[1].encode("UTF-8")).digest() # Extract and hash password from tuple
        password = self.hex(base64.b64encode(password)) # Encode password

        school   = credentials[2]
        brin     = credentials[3]

        loginURL = "http://somtoday.nl/" + school + "/services/mobile/v10/Login/CheckMultiLoginB64/" +  username + "/" + password + "/" + brin

        response = self.getJSON(loginURL) # Perform request

        if response["error"] == "SUCCESS":
            self.leerlingID = str(response["leerlingen"][0]["leerlingId"])
            self.username = username
            self.password = password
            self.school = school
            self.brin = brin
            self.baseURL = "http://somtoday.nl/" + school + "/services/mobile/v10/"
        elif response["error"] == "FEATURE_NOT_ACTIVATED":
            raise FEATURE_NOT_ACTIVATED("Your school isn't supported by SOMToday")
        elif response["error"] == "FAILED_AUTHENTICATION":
            raise FAILED_AUTENTICATION("Your credentials aren't right")
        elif response["error"] == "FAILED_OTHER_TYPE":
            raise FAILED_OTHER_TYPE("Your type of account isn't supported")
        else:
            raise SOMTODAY_ERROR_UNKNOWN("An unknown error has occured")

    def getTimetable(self, days = 0):
        days = int(time.time() + (days * 86400)) * 1000

        URL = self.baseURL + "Agenda/GetMultiStudentAgendaB64/" + self.username + "/" + self.password + "/" + self.brin + "/" + str(days) + "/" + self.leerlingID
        timetable = self.getJSON(URL)["data"]

        return timetable

    def getHW(self, daysMax = 31):
        URL = self.baseURL + "Agenda/GetMultiStudentAgendaHuiswerkMetMaxB64/" + self.username + "/" + self.password + "/" + self.brin + "/" + str(daysMax) + "/" + self.leerlingID
        hw = self.getJSON(URL)["data"]

        return hw

    def getGrades(self, mode = 1):

        if mode == 0:
            URL = self.baseURL + "Cijfers/GetMultiCijfersRecentB64/" + self.username + "/" + self.password + "/" + self.brin + "/" + self.leerlingID
        elif mode == 1:
            URL = self.baseURL + "Cijfers/GetCijfersRecentMetMaxB64/" + self.username + "/" + self.password + "/" + self.brin + "/1000"

        grades = self.getJSON(URL)["data"]
        return grades

    def changeHwDone(self, appointmentID, hwID, done = True):

        URL = self.baseURL + "Agenda/HuiswerkDoneB64/" + self.username + "/" + self.password + "/" + self.brin + "/" + str(appointmentID) + "/" + str(hwID) + "/"

        if done:
            URL += "1"
        else:
            URL += "0"

        response = self.getJSON(URL)

        if response["status"] == "OK":
            return True
        else:
            return False

    @staticmethod
    def hex(toHex):
        return str(codecs.encode(toHex, "hex").decode())

    @staticmethod
    def getJSON(URL):
        
        r = requests.get(URL)
        response = r.json()

        return response