Radarbilder levereras i 8bitars GeoTIFF-format.  
 
Bilderna är i SWEREF99-projektion med hörnkoordinaterna:
nedre vänstra: 126648 5983984 
övre högra: 1075693 7771252
 
Filerna med data innehåller värden 0-255.
0 = No echo. Mätpunkten är innanför radarns täckningsområde men inget eko är rapporterat.
255= No data. Mätpunkten är utanför radarns täckningsområde.
Övriga värden i filen kan översättas till dBZ genom att applicera en offset och en gain.
I de här filerna gäller att offset = -30 och gain = 0.4
Därmed är värdet i dBz = pixelvalue*gain + offset

Filen palett.txt kan användas i QGIS för att få GeoTIFF-bilden att se ut som motsvarande radarbild i PNG-format.
