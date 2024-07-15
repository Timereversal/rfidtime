# rfidtime
####  chip type Alien H3 9654 
# v1
# PC - Decoder/Encoder - 1 Antenna.
# Almacenar los valores por cada EPC - RSSI- Time [Database or something else]

# flow 
# 0x21 get reader information
# 0x51  stop fast inventory
# 0x7f load/modify reader profile

//
# start inventory
# Data: 0c009a010300010000000032ef  (select command) (send like 3 times)
# Data: 09000104fe00803280be
# Data: 05007fcd9303
# Data: 09000104fe00803280be
# Data: 0700010101001e4b


# Data: 09000104fe0180325ce4
# anatomy Data: 09 00 01 04(Qvalue) fe(Session) 01(MaskMem) 8032(MaskAddres) 5c(MaskLen) e4



#########################################
#########################################
TODO List
 -  Define clear Inventory initialization Pro and cons.
 -  Find a way to send a command to stop realtime inventory on demand like graceful shutdown
 -  Test Realtime inventory with more labels.
 -  Find a way to perform Answer Mode inventory/ equipment no anwser always to the same command.
 -  add file log.
 -  Define Group structure for like RunnerID , Group/Category.
   - runner range [1-10000] 2 byte
   - runner category [1-16] 1 byte
        - Juvenil/Libre/Master/ etc 
        - Subcategoria(Edades)
        - Sexo [1bit]
        - Distance [8-max]
   - experimental/reserved  2 bytes

Concurrency :
Notes :
Go Routines:
 - Just like any other function, you can pass it parameters to initialize its state.However
    any values returned by the function are ignored (bodner)
 - 
